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

	lp := make([]string, 0, len(importPaths)+1)
	lp = append(lp, filepath.Dir(abs))
	lp = append(lp, importPaths...)

	out := &Loaded{
		Files: map[string]*File{},
	}

	visiting := map[string]bool{}
	root, err := loadOne(abs, lp, visiting, out.Files)
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
			for _, e := range s.Elements {
				r, ok := e.(*proto.RPC)
				if !ok {
					continue
				}
				rpc := Rpc{
					Name:     r.Name,
					Request:  r.RequestType,
					Response: r.ReturnsType,
				}
				if r.Comment != nil {
					rpc.Http = parseHTTPRule(r.Comment.Message())
				}
				svc.Rpcs = append(svc.Rpcs, rpc)
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

