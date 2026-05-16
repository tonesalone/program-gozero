package protosrc

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/emicklei/proto"
)

var httpRe = regexp.MustCompile(`(?i)\bhttp\s+([A-Z]+)\s+(\S+)`)

func Load(path string, importPaths []string) (*Loaded, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	lp := make([]string, 0, len(importPaths)+2)
	lp = append(lp, filepath.Dir(abs))
	if wd, err := os.Getwd(); err == nil {
		lp = append(lp, wd)
	}
	for _, p := range importPaths {
		if p == "" {
			continue
		}
		ap, err := filepath.Abs(p)
		if err == nil {
			p = ap
		}
		lp = append(lp, p)
	}
	seen := map[string]bool{}
	dedup := make([]string, 0, len(lp))
	for _, p := range lp {
		if p == "" {
			continue
		}
		p = filepath.Clean(p)
		if seen[p] {
			continue
		}
		seen[p] = true
		dedup = append(dedup, p)
	}

	out := &Loaded{
		Files: map[string]*File{},
	}

	visiting := map[string]bool{}
	root, err := loadOne(abs, dedup, visiting, out.Files)
	if err != nil {
		return nil, err
	}
	out.Root = root
	return out, nil
}

func loadOne(absPath string, importPaths []string, visiting map[string]bool, files map[string]*File) (*File, error) {
	if visiting[absPath] {
		return nil, fmt.Errorf("import cycle: %s", absPath)
	}
	if f, ok := files[absPath]; ok {
		return f, nil
	}
	visiting[absPath] = true
	defer delete(visiting, absPath)

	fh, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	parser := proto.NewParser(bufio.NewReader(fh))
	def, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	f := &File{Path: absPath}

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) {
			f.Package = p.Name
		}),
		proto.WithImport(func(i *proto.Import) {
			f.Imports = append(f.Imports, strings.Trim(i.Filename, `"`))
		}),
		proto.WithEnum(func(e *proto.Enum) {
			en := Enum{Name: e.Name}
			for _, el := range e.Elements {
				v, ok := el.(*proto.EnumField)
				if !ok {
					continue
				}
				en.Values = append(en.Values, EnumValue{
					Name:  v.Name,
					Index: v.Integer,
				})
			}
			f.Enums = append(f.Enums, en)
		}),
		proto.WithService(func(s *proto.Service) {
			svc := Service{Name: s.Name}
			var pending []string
			flushPending := func() string {
				if len(pending) == 0 {
					return ""
				}
				s := strings.Join(pending, "\n")
				pending = pending[:0]
				return s
			}
			for _, e := range s.Elements {
				switch el := e.(type) {
				case *proto.Comment:
					pending = append(pending, el.Message())
					continue
				case *proto.RPC:
					r := el
					rpc := Rpc{
						Name:     r.Name,
						Request:  r.RequestType,
						Response: r.ReturnsType,
					}
					var cands []string
					if r.Comment != nil {
						cands = append(cands, r.Comment.Message())
					}
					if r.InlineComment != nil {
						cands = append(cands, r.InlineComment.Message())
					}
					if s := flushPending(); s != "" {
						cands = append(cands, s)
					}
					for _, c := range cands {
						rule := parseHTTPRule(c)
						if rule.Method != "" && rule.Path != "" {
							rpc.Http = rule
							break
						}
					}
					svc.Rpcs = append(svc.Rpcs, rpc)
				default:
					pending = pending[:0]
				}
			}
			f.Services = append(f.Services, svc)
		}),
	)

	collectMessages(def.Elements, nil, &f.Messages)

	files[absPath] = f

	for _, imp := range f.Imports {
		impAbs, err := resolveImport(absPath, imp, importPaths)
		if err != nil {
			return nil, err
		}
		if _, err := loadOne(impAbs, importPaths, visiting, files); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func parseHTTPRule(comment string) HttpRule {
	m := httpRe.FindStringSubmatch(comment)
	if len(m) != 3 {
		return HttpRule{}
	}
	return HttpRule{
		Method: strings.ToUpper(m[1]),
		Path:   m[2],
	}
}

func resolveImport(fromAbs string, importName string, importPaths []string) (string, error) {
	if filepath.IsAbs(importName) {
		return importName, nil
	}

	cand := filepath.Join(filepath.Dir(fromAbs), importName)
	if st, err := os.Stat(cand); err == nil && !st.IsDir() {
		return filepath.Abs(cand)
	}

	for _, base := range importPaths {
		cand := filepath.Join(base, importName)
		if st, err := os.Stat(cand); err == nil && !st.IsDir() {
			return filepath.Abs(cand)
		}
	}

	return "", errors.New("cannot resolve import " + importName + " from " + fromAbs)
}

func collectMessages(elements []proto.Visitee, prefix []string, out *[]Message) {
	for _, e := range elements {
		m, ok := e.(*proto.Message)
		if !ok {
			continue
		}
		nameParts := make([]string, 0, len(prefix)+1)
		nameParts = append(nameParts, prefix...)
		nameParts = append(nameParts, m.Name)
		msg := Message{Name: strings.Join(nameParts, ".")}
		for _, me := range m.Elements {
			switch el := me.(type) {
			case *proto.NormalField:
				msg.Fields = append(msg.Fields, Field{
					Name:     el.Name,
					TypeName: el.Type,
					Repeated: el.Repeated,
				})
			case *proto.MapField:
				msg.Fields = append(msg.Fields, Field{
					Name:     el.Name,
					TypeName: "map",
					MapKey:   el.KeyType,
					MapValue: el.Type,
				})
			}
		}
		*out = append(*out, msg)

		nextPrefix := append(nameParts[:0:0], nameParts...)
		collectMessages(m.Elements, nextPrefix, out)
	}
}
