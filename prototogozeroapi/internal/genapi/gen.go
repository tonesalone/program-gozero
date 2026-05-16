package genapi

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"prototogozeroapi/internal/protosrc"
)

type msgInfo struct {
	filePath string
	pkg      string
	msg      protosrc.Message
}

func Generate(loaded *protosrc.Loaded) (string, error) {
	if loaded == nil || loaded.Root == nil {
		return "", fmt.Errorf("nil input")
	}

	root := loaded.Root

	msgs := map[string]msgInfo{}
	enums := map[string]bool{}
	for _, f := range loaded.Files {
		for _, m := range f.Messages {
			msgs[fullName(f.Package, m.Name)] = msgInfo{
				filePath: f.Path,
				pkg:      f.Package,
				msg:      m,
			}
		}
		for _, e := range f.Enums {
			enums[fullName(f.Package, e.Name)] = true
		}
	}

	used := map[string]bool{}
	queue := make([]string, 0, 32)

	for _, m := range root.Messages {
		fn := fullName(root.Package, m.Name)
		if !used[fn] {
			used[fn] = true
			queue = append(queue, fn)
		}
	}

	for _, s := range root.Services {
		for _, r := range s.Rpcs {
			for _, tn := range []string{r.Request, r.Response} {
				fn, ok, err := resolveTypeRef(root.Package, "", tn, func(fn string) bool {
					return enums[fn] || msgs[fn].msg.Name != ""
				})
				if err != nil {
					return "", err
				}
				if !ok {
					continue
				}
				if enums[fn] {
					continue
				}
				if _, exists := msgs[fn]; !exists {
					return "", fmt.Errorf("unknown type: %s", tn)
				}
				if !used[fn] {
					used[fn] = true
					queue = append(queue, fn)
				}
			}
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		mi, ok := msgs[cur]
		if !ok {
			return "", fmt.Errorf("unknown message: %s", cur)
		}
		for _, f := range mi.msg.Fields {
			refs := fieldTypeRefs(f)
			for _, r := range refs {
				fn, ok, err := resolveTypeRef(mi.pkg, mi.msg.Name, r, func(fn string) bool {
					return enums[fn] || msgs[fn].msg.Name != ""
				})
				if err != nil {
					return "", err
				}
				if !ok {
					continue
				}
				if enums[fn] {
					continue
				}
				if _, exists := msgs[fn]; !exists {
					return "", fmt.Errorf("unknown type: %s", r)
				}
				if !used[fn] {
					used[fn] = true
					queue = append(queue, fn)
				}
			}
		}
	}

	nameMap := buildNameMap(used, msgs)

	var b strings.Builder
	b.WriteString(`syntax = "v1"` + "\n\n")

	b.WriteString("type (\n")
	for _, m := range root.Messages {
		fn := fullName(root.Package, m.Name)
		apiName := nameMap[fn]
		b.WriteString("\n    " + apiName + " {\n")
		if err := writeFields(&b, root.Package, m.Fields, nameMap, enums); err != nil {
			return "", err
		}
		b.WriteString("    }\n")
	}

	extras := make([]string, 0, len(used))
	for fn := range used {
		if !isFromRoot(root.Package, fn, root.Messages) {
			extras = append(extras, fn)
		}
	}
	sort.Strings(extras)
	for _, fn := range extras {
		mi := msgs[fn]
		apiName := nameMap[fn]
		b.WriteString("\n    " + apiName + " {\n")
		if err := writeFields(&b, mi.pkg, mi.msg.Fields, nameMap, enums); err != nil {
			return "", err
		}
		b.WriteString("    }\n")
	}

	b.WriteString("\n)\n\n")

	for _, s := range root.Services {
		b.WriteString(fmt.Sprintf("service %s {\n\n", strings.ToLower(s.Name)+"-api"))
		for _, r := range s.Rpcs {
			method, path := routeFor(root.Package, s.Name, r)
			b.WriteString(fmt.Sprintf("    @handler %sHandler\n", r.Name))
			req, err := apiTypeName(root.Package, r.Request, nameMap, enums)
			if err != nil {
				return "", err
			}
			resp, err := apiTypeName(root.Package, r.Response, nameMap, enums)
			if err != nil {
				return "", err
			}
			b.WriteString(fmt.Sprintf("    %s %s (%s) returns (%s)\n\n",
				strings.ToLower(method), path, req, resp))
		}
		b.WriteString("}\n")
	}

	return b.String(), nil
}

func isFromRoot(rootPkg string, fn string, rootMsgs []protosrc.Message) bool {
	for _, m := range rootMsgs {
		if fullName(rootPkg, m.Name) == fn {
			return true
		}
	}
	return false
}

func buildNameMap(used map[string]bool, msgs map[string]msgInfo) map[string]string {
	baseCount := map[string]int{}
	for fn := range used {
		mi := msgs[fn]
		baseCount[apiBaseName(mi.msg.Name)]++
	}

	nameMap := map[string]string{}
	for fn := range used {
		mi := msgs[fn]
		base := apiBaseName(mi.msg.Name)
		if baseCount[base] == 1 {
			nameMap[fn] = base
			continue
		}
		prefix := strings.ReplaceAll(mi.pkg, ".", "_")
		if prefix == "" {
			base := strings.TrimSuffix(filepath.Base(mi.filePath), filepath.Ext(mi.filePath))
			prefix = base
		}
		nameMap[fn] = prefix + "_" + base
	}
	return nameMap
}

func writeFields(b *strings.Builder, rootPkg string, fields []protosrc.Field, nameMap map[string]string, enums map[string]bool) error {
	for _, fld := range fields {
		tp, err := apiFieldType(rootPkg, fld, nameMap, enums)
		if err != nil {
			return err
		}
		tags := fmt.Sprintf("`json:\"%s,optional\" form:\"%s,optional\"`", fld.Name, fld.Name)
		b.WriteString(fmt.Sprintf("        %s %s %s\n", fld.Name, tp, tags))
	}
	return nil
}

func apiFieldType(rootPkg string, fld protosrc.Field, nameMap map[string]string, enums map[string]bool) (string, error) {
	if fld.MapKey != "" || fld.MapValue != "" {
		k, err := apiScalarType(rootPkg, fld.MapKey, nameMap, enums)
		if err != nil {
			return "", err
		}
		v, err := apiScalarType(rootPkg, fld.MapValue, nameMap, enums)
		if err != nil {
			return "", err
		}
		return "map[" + k + "]" + v, nil
	}

	tp, err := apiScalarType(rootPkg, fld.TypeName, nameMap, enums)
	if err != nil {
		return "", err
	}
	if fld.Repeated {
		return "[]" + tp, nil
	}
	return tp, nil
}

func apiScalarType(rootPkg string, typeName string, nameMap map[string]string, enums map[string]bool) (string, error) {
	if typeName == "" {
		return "", fmt.Errorf("empty type")
	}

	if tp, ok := scalarMap(typeName); ok {
		return tp, nil
	}

	fn, ok, err := resolveTypeRef(rootPkg, "", typeName, func(fn string) bool {
		return enums[fn] || nameMap[fn] != ""
	})
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("unsupported type: %s", typeName)
	}
	if enums[fn] {
		return "int32", nil
	}
	if tp, ok := nameMap[fn]; ok && tp != "" {
		return tp, nil
	}
	return "", fmt.Errorf("unknown type: %s", typeName)
}

func apiTypeName(rootPkg string, typeName string, nameMap map[string]string, enums map[string]bool) (string, error) {
	if tp, ok := scalarMap(typeName); ok {
		return tp, nil
	}
	fn, ok, err := resolveTypeRef(rootPkg, "", typeName, func(fn string) bool {
		return enums[fn] || nameMap[fn] != ""
	})
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("unsupported type: %s", typeName)
	}
	if enums[fn] {
		return "int32", nil
	}
	if tp, ok := nameMap[fn]; ok && tp != "" {
		return tp, nil
	}
	return "", fmt.Errorf("unknown type: %s", typeName)
}

func routeFor(pkg string, svc string, r protosrc.Rpc) (string, string) {
	if r.Http.Method != "" && r.Http.Path != "" {
		return r.Http.Method, r.Http.Path
	}
	pkgPath := strings.ReplaceAll(pkg, ".", "/")
	pkgPath = strings.ToLower(pkgPath)
	svcSeg := strings.ToLower(svc)
	rpcSeg := strings.ToLower(r.Name)
	var parts []string
	parts = append(parts, "")
	if pkgPath != "" {
		parts = append(parts, pkgPath)
	}
	parts = append(parts, svcSeg, rpcSeg)
	return "POST", strings.Join(parts, "/")
}

func fieldTypeRefs(f protosrc.Field) []string {
	out := make([]string, 0, 2)
	if f.MapKey != "" {
		out = append(out, f.MapKey)
	}
	if f.MapValue != "" {
		out = append(out, f.MapValue)
	}
	if f.MapKey == "" && f.MapValue == "" {
		out = append(out, f.TypeName)
	}
	return out
}

func scalarMap(tp string) (string, bool) {
	switch tp {
	case "string":
		return "string", true
	case "bool":
		return "bool", true
	case "bytes":
		return "string", true
	case "int32", "sint32", "sfixed32":
		return "int32", true
	case "uint32", "fixed32":
		return "uint32", true
	case "int64", "sint64", "sfixed64":
		return "int64", true
	case "uint64", "fixed64":
		return "uint64", true
	case "float":
		return "float32", true
	case "double":
		return "float64", true
	}
	if tp == "google.protobuf.Timestamp" || tp == ".google.protobuf.Timestamp" {
		return "string", true
	}
	return "", false
}

func resolveTypeRef(pkg string, scope string, typeName string, exists func(string) bool) (string, bool, error) {
	if typeName == "" {
		return "", false, nil
	}
	if _, ok := scalarMap(typeName); ok {
		return "", false, nil
	}

	if strings.HasPrefix(typeName, ".") {
		return strings.TrimPrefix(typeName, "."), true, nil
	}

	tn := typeName
	if strings.Contains(tn, ".") {
		if pkg != "" {
			cand := pkg + "." + tn
			if exists(cand) {
				return cand, true, nil
			}
		}
		if exists(tn) {
			return tn, true, nil
		}
		if pkg == "" {
			return tn, true, nil
		}
		return pkg + "." + tn, true, nil
	}

	if scope != "" {
		parts := strings.Split(scope, ".")
		for i := len(parts); i >= 1; i-- {
			pref := strings.Join(parts[:i], ".")
			var cand string
			if pkg != "" {
				cand = pkg + "." + pref + "." + tn
			} else {
				cand = pref + "." + tn
			}
			if exists(cand) {
				return cand, true, nil
			}
		}
	}

	return fullName(pkg, tn), true, nil
}

func fullName(pkg string, name string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

func apiBaseName(protoName string) string {
	return strings.ReplaceAll(protoName, ".", "_")
}

