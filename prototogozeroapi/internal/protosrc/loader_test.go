package protosrc

import (
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

