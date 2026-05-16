package protosrc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ParsesRpcHttpComment(t *testing.T) {
	root := filepath.Join("..", "testdata", "proto", "simple")
	protoPath := filepath.Join(root, "simple.proto")

	loaded, err := Load(protoPath, []string{root})
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	f := loaded.Root

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

func TestLoad_ResolvesImport_FromWorkingDirAndImporterDir(t *testing.T) {
	root := filepath.Join("..", "testdata", "proto", "multi")
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs root: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(rootAbs); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	loaded, err := Load(filepath.Join("imserver", "imserver.proto"), nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	want := []string{
		filepath.Join(rootAbs, "imserver", "imserver.proto"),
		filepath.Join(rootAbs, "public", "public.proto"),
		filepath.Join(rootAbs, "types", "value.proto"),
		filepath.Join(rootAbs, "types", "nested.proto"),
	}
	for _, p := range want {
		p = filepath.Clean(p)
		if _, ok := loaded.Files[p]; !ok {
			t.Fatalf("missing loaded file: %s", p)
		}
	}
}

func TestLoad_ParsesRpcHttpComment_FromLooseOrInlineComments(t *testing.T) {
	protoPath := filepath.Join("..", "testdata", "proto", "httpcomment", "httpcomment.proto")

	loaded, err := Load(protoPath, nil)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Root.Services) != 1 {
		t.Fatalf("unexpected services")
	}
	rpcs := loaded.Root.Services[0].Rpcs
	if len(rpcs) != 2 {
		t.Fatalf("unexpected rpcs")
	}

	if rpcs[0].Http.Path != "/cs/v1/b/imserver/dialog_route/delete" {
		t.Fatalf("path0 = %q", rpcs[0].Http.Path)
	}
	if rpcs[1].Http.Path != "/cs/v1/b/imserver/dialog_route/inline" {
		t.Fatalf("path1 = %q", rpcs[1].Http.Path)
	}
}
