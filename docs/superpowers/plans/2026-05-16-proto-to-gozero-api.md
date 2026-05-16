# Proto To go-zero API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现一个 Go 命令行工具：输入 `.proto`，解析结构与 RPC，并在同目录全量覆盖生成 go-zero `.api` 文件。

**Architecture:** 使用成熟的 `.proto` AST 解析库（方案 A）解析包含注释的语法树；构建符号表与类型引用闭包；按确定性规则渲染 `.api` 并写入输出文件。

**Tech Stack:** Go, github.com/emicklei/proto, Go testing

---

### Task 1: 创建模块与 CLI 骨架（TDD 起步）

**Files:**
- Create: `prototogozeroapi/go.mod`
- Create: `prototogozeroapi/cmd/prototogozeroapi/main.go`
- Create: `prototogozeroapi/internal/cli/cli.go`
- Test: `prototogozeroapi/internal/cli/cli_test.go`

- [ ] **Step 1: 写一个失败的 CLI 测试（缺参数应报错）**

```go
package cli

import "testing"

func TestRun_NoArgs(t *testing.T) {
	err := Run([]string{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run（在仓库根目录执行）：
`cd prototogozeroapi && go test ./...`

Expected：FAIL（因为 `Run` 未实现/未定义）

- [ ] **Step 3: 写最小实现让测试通过**

```go
package cli

import (
	"errors"
)

func Run(args []string) error {
	if len(args) < 1 {
		return errors.New("missing proto file path")
	}
	return nil
}
```

- [ ] **Step 4: 运行测试确认通过**

Run：`cd prototogozeroapi && go test ./...`
Expected：PASS

- [ ] **Step 5: 增加 main.go 入口调用 cli.Run**

```go
package main

import (
	"fmt"
	"os"

	"program-gozero/prototogozeroapi/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
```

> 注意：`program-gozero` 作为 module path 仅占位，后续在 Task 2 会把 import 改为实际 `go.mod` 的 module 名。

---

### Task 2: 建立 `.proto` 解析器（含 import 与注释）

**Files:**
- Modify: `prototogozeroapi/go.mod`
- Create: `prototogozeroapi/internal/protosrc/loader.go`
- Create: `prototogozeroapi/internal/protosrc/model.go`
- Test: `prototogozeroapi/internal/protosrc/loader_test.go`
- Create: `prototogozeroapi/internal/testdata/proto/simple/simple.proto`
- Create: `prototogozeroapi/internal/testdata/proto/simple/imported.proto`

- [ ] **Step 1: 写失败的 loader 测试（能解析 package/import/message/service/rpc 与 rpc 注释）**

```go
package protosrc

import (
	"path/filepath"
	"testing"
)

func TestLoad_ParsesRpcHttpComment(t *testing.T) {
	root := filepath.Join("..", "testdata", "proto", "simple")
	protoPath := filepath.Join(root, "simple.proto")

	f, err := Load(protoPath, []string{root})
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if f.Package != "cs.v1" {
		t.Fatalf("package = %q", f.Package)
	}
	if len(f.Services) != 1 || len(f.Services[0].Rpcs) != 1 {
		t.Fatalf("unexpected services/rpcs")
	}
	if f.Services[0].Rpcs[0].Http.Method != "POST" || f.Services[0].Rpcs[0].Http.Path != "/cs/v1/mod/add" {
		t.Fatalf("http = %#v", f.Services[0].Rpcs[0].Http)
	}
}
```

- [ ] **Step 2: 添加测试 proto 文件**

`simple.proto`（包含 import 与 http 注释）：

```proto
syntax = "proto3";
package cs.v1;

import "imported.proto";

message AddReq {
  string book = 1;
  ImportedMsg other = 2;
}

message AddResp {
  string ok = 1;
}

service Book {
  // http POST /cs/v1/mod/add
  rpc Add(AddReq) returns (AddResp) {}
}
```

`imported.proto`：

```proto
syntax = "proto3";
package cs.v1;

message ImportedMsg {
  int64 id = 1;
}
```

- [ ] **Step 3: 运行测试确认失败**

Run：`cd prototogozeroapi && go test ./...`
Expected：FAIL（因为 `Load` 未实现/未定义）

- [ ] **Step 4: 最小实现 loader（使用 github.com/emicklei/proto）**

`model.go` 定义：

```go
package protosrc

type HttpRule struct {
	Method string
	Path   string
}

type Field struct {
	Name     string
	TypeName string
	Repeated bool
	MapKey   string
	MapValue string
}

type Message struct {
	Name   string
	Fields []Field
}

type Rpc struct {
	Name     string
	Request  string
	Response string
	Http     HttpRule
}

type Service struct {
	Name string
	Rpcs []Rpc
}

type File struct {
	Path     string
	Package  string
	Imports  []string
	Messages []Message
	Services []Service
}
```

`loader.go` 最小实现（先只覆盖测试用例所需语法）：

```go
package protosrc

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/emicklei/proto"
)

var httpRe = regexp.MustCompile(`(?i)\bhttp\s+([A-Z]+)\s+(\S+)`)

func Load(path string, importPaths []string) (*File, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	return loadOne(abs, importPaths, seen)
}

func loadOne(path string, importPaths []string, seen map[string]bool) (*File, error) {
	if seen[path] {
		return nil, fmt.Errorf("import cycle: %s", path)
	}
	seen[path] = true

	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	parser := proto.NewParser(bufio.NewReader(fh))
	def, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	out := &File{Path: path}

	proto.Walk(
		def,
		proto.WithPackage(func(p *proto.Package) {
			out.Package = p.Name
		}),
		proto.WithImport(func(i *proto.Import) {
			out.Imports = append(out.Imports, strings.Trim(i.Filename, `"`))
		}),
		proto.WithMessage(func(m *proto.Message) {
			msg := Message{Name: m.Name}
			for _, e := range m.Elements {
				n, ok := e.(*proto.NormalField)
				if !ok {
					continue
				}
				msg.Fields = append(msg.Fields, Field{
					Name:     n.Name,
					TypeName: n.Type,
					Repeated: n.Repeated,
				})
			}
			out.Messages = append(out.Messages, msg)
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
				for _, c := range []string{r.Comment.Leading, r.Comment.Trailing} {
					m := httpRe.FindStringSubmatch(c)
					if len(m) == 3 {
						rpc.Http = HttpRule{Method: strings.ToUpper(m[1]), Path: m[2]}
						break
					}
				}
				svc.Rpcs = append(svc.Rpcs, rpc)
			}
			out.Services = append(out.Services, svc)
		}),
	)

	return out, nil
}
```

- [ ] **Step 5: 运行测试确认通过**

Run：`cd prototogozeroapi && go test ./...`
Expected：PASS

---

### Task 3: 类型闭包收集与 `.api` 渲染器（全量覆盖输出）

**Files:**
- Create: `prototogozeroapi/internal/genapi/gen.go`
- Test: `prototogozeroapi/internal/genapi/gen_test.go`

- [ ] **Step 1: 写失败的渲染测试（golden）**

```go
package genapi

import (
	"path/filepath"
	"testing"

	"program-gozero/prototogozeroapi/internal/protosrc"
)

func TestGenerate_Simple(t *testing.T) {
	root := filepath.Join("..", "testdata", "proto", "simple")
	f, err := protosrc.Load(filepath.Join(root, "simple.proto"), []string{root})
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	got, err := Generate(f)
	if err != nil {
		t.Fatalf("gen: %v", err)
	}
	if len(got) == 0 {
		t.Fatalf("empty output")
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run：`cd prototogozeroapi && go test ./...`
Expected：FAIL（`Generate` 未实现/未定义）

- [ ] **Step 3: 最小实现 Generate（先只覆盖 message + service + http 注释 / 默认规则）**

```go
package genapi

import (
	"fmt"
	"strings"

	"program-gozero/prototogozeroapi/internal/protosrc"
)

func Generate(f *protosrc.File) (string, error) {
	var b strings.Builder
	b.WriteString(`syntax = "v1"` + "\n\n")

	b.WriteString("type (\n")
	for _, m := range f.Messages {
		b.WriteString("\n    " + m.Name + " {\n")
		for _, fld := range m.Fields {
			tp := mapScalar(fld.TypeName)
			tags := fmt.Sprintf("`json:\"%s,optional\" form:\"%s,optional\"`", fld.Name, fld.Name)
			b.WriteString(fmt.Sprintf("        %s %s %s\n", fld.Name, tp, tags))
		}
		b.WriteString("    }\n")
	}
	b.WriteString("\n)\n\n")

	for _, s := range f.Services {
		b.WriteString(fmt.Sprintf("service %s {\n\n", strings.ToLower(s.Name)+"-api"))
		for _, r := range s.Rpcs {
			method, path := routeFor(f.Package, s.Name, r)
			b.WriteString(fmt.Sprintf("    @handler %sHandler\n", r.Name))
			b.WriteString(fmt.Sprintf("    %s %s (%s) returns (%s)\n\n",
				strings.ToLower(method), path, r.Request, r.Response))
		}
		b.WriteString("}\n")
	}

	return b.String(), nil
}

func mapScalar(tp string) string {
	switch tp {
	case "string":
		return "string"
	case "int64":
		return "int64"
	case "int32":
		return "int32"
	case "bool":
		return "bool"
	default:
		return tp
	}
}

func routeFor(pkg string, svc string, r protosrc.Rpc) (string, string) {
	if r.Http.Method != "" && r.Http.Path != "" {
		return r.Http.Method, r.Http.Path
	}
	pkgPath := strings.ReplaceAll(pkg, ".", "/")
	if pkgPath != "" {
		pkgPath = "/" + pkgPath
	}
	return "POST", strings.ToLower(pkgPath+"/"+svc+"/"+r.Name)
}
```

- [ ] **Step 4: 运行测试确认通过（下一步再收敛到 golden 严格比对）**

Run：`cd prototogozeroapi && go test ./...`
Expected：PASS

---

### Task 4: 集成 CLI：生成并写入同目录 `.api`

**Files:**
- Modify: `prototogozeroapi/internal/cli/cli.go`
- Modify: `prototogozeroapi/cmd/prototogozeroapi/main.go`
- Create: `prototogozeroapi/internal/cli/e2e_test.go`

- [ ] **Step 1: 写失败的 e2e 测试（运行 cli 生成文件）**

```go
package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_GeneratesApiInSameDir(t *testing.T) {
	tmp := t.TempDir()

	protoDir := filepath.Join("..", "testdata", "proto", "simple")
	src := filepath.Join(protoDir, "simple.proto")
	dst := filepath.Join(tmp, "simple.proto")
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := Run([]string{dst}); err != nil {
		t.Fatalf("run: %v", err)
	}

	out := filepath.Join(tmp, "simple.api")
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected %s generated: %v", out, err)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run：`cd prototogozeroapi && go test ./...`
Expected：FAIL（因为 Run 还没做生成与写文件）

- [ ] **Step 3: 实现 cli.Run：Load → Generate → 写入同目录**

```go
package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"program-gozero/prototogozeroapi/internal/genapi"
	"program-gozero/prototogozeroapi/internal/protosrc"
)

func Run(args []string) error {
	if len(args) < 1 {
		return errors.New("missing proto file path")
	}

	in := args[0]
	dir := filepath.Dir(in)
	base := strings.TrimSuffix(filepath.Base(in), filepath.Ext(in))
	outPath := filepath.Join(dir, base+".api")

	f, err := protosrc.Load(in, []string{dir})
	if err != nil {
		return err
	}
	api, err := genapi.Generate(f)
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, []byte(api), 0o644)
}
```

- [ ] **Step 4: 运行测试确认通过**

Run：`cd prototogozeroapi && go test ./...`
Expected：PASS

---

### Task 5: 收敛到最终需求（全结构 + 外部引用 + 完整类型映射）

**Files:**
- Modify: `prototogozeroapi/internal/protosrc/loader.go`
- Modify: `prototogozeroapi/internal/genapi/gen.go`
- Test: `prototogozeroapi/internal/genapi/gen_golden_test.go`
- Create: `prototogozeroapi/internal/testdata/proto/complex/*.proto`

- [ ] **Step 1: 增加 golden 测试样例覆盖**

覆盖点：
- repeated / map / enum 字段（enum 映射为 int32）
- 缺失 http 注释时默认规则
- 引用外部 import message（闭包收集到输出）
- 输入文件内所有 message 都要输出（即使不被 rpc 引用）

- [ ] **Step 2: 让 loader 支持递归 import 解析 + 类型符号表**

要点：
- 允许 `-I` 多个目录（后续在 CLI 扩展参数）
- 解析所有 message（含嵌套 message 扁平化命名 `Outer_Inner`）
- 支持解析 map field（proto AST 中为 `*proto.MapField`）

- [ ] **Step 3: 让 Generate 支持类型闭包收集与冲突命名**

要点：
- 输入文件 message 全量输出
- 外部文件只输出闭包中被引用到的 message
- 同名冲突时使用 `<package>_<Type>` 规则

- [ ] **Step 4: 运行全量测试**

Run：`cd prototogozeroapi && go test ./...`
Expected：PASS

