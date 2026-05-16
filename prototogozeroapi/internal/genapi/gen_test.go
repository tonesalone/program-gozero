package genapi

import (
	"path/filepath"
	"strings"
	"testing"

	"prototogozeroapi/internal/protosrc"
)

func TestGenerate_Simple(t *testing.T) {
	root := filepath.Join("..", "testdata", "proto", "simple")
	loaded, err := protosrc.Load(filepath.Join(root, "simple.proto"), []string{root})
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	got, err := Generate(loaded)
	if err != nil {
		t.Fatalf("gen: %v", err)
	}
	if !strings.Contains(got, `syntax = "v1"`) {
		t.Fatalf("missing syntax")
	}
	if !strings.Contains(got, "service book-api") {
		t.Fatalf("missing service")
	}
	if !strings.Contains(got, "post /cs/v1/mod/add (AddReq) returns (AddResp)") {
		t.Fatalf("missing route")
	}
	if !strings.Contains(got, "ImportedMsg") {
		t.Fatalf("missing imported type")
	}
}

