package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_GeneratesApiInSameDir(t *testing.T) {
	tmp := t.TempDir()

	protoDir := filepath.Join("..", "testdata", "proto", "simple")
	for _, name := range []string{"simple.proto", "imported.proto"} {
		src := filepath.Join(protoDir, name)
		dst := filepath.Join(tmp, name)
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if err := os.WriteFile(dst, b, 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	if err := Run([]string{filepath.Join(tmp, "simple.proto")}); err != nil {
		t.Fatalf("run: %v", err)
	}

	out := filepath.Join(tmp, "simple.api")
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected %s generated: %v", out, err)
	}
}

func TestRun_IncludePath(t *testing.T) {
	tmp := t.TempDir()
	a := filepath.Join(tmp, "a")
	b := filepath.Join(tmp, "b")
	if err := os.MkdirAll(a, 0o755); err != nil {
		t.Fatalf("mkdir a: %v", err)
	}
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatalf("mkdir b: %v", err)
	}

	protoDir := filepath.Join("..", "testdata", "proto", "simple")
	simpleSrc := filepath.Join(protoDir, "simple.proto")
	importedSrc := filepath.Join(protoDir, "imported.proto")

	simpleDst := filepath.Join(a, "simple.proto")
	importedDst := filepath.Join(b, "imported.proto")

	b1, err := os.ReadFile(simpleSrc)
	if err != nil {
		t.Fatalf("read simple: %v", err)
	}
	if err := os.WriteFile(simpleDst, b1, 0o644); err != nil {
		t.Fatalf("write simple: %v", err)
	}
	b2, err := os.ReadFile(importedSrc)
	if err != nil {
		t.Fatalf("read imported: %v", err)
	}
	if err := os.WriteFile(importedDst, b2, 0o644); err != nil {
		t.Fatalf("write imported: %v", err)
	}

	if err := Run([]string{"-I", b, simpleDst}); err != nil {
		t.Fatalf("run: %v", err)
	}

	out := filepath.Join(a, "simple.api")
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected %s generated: %v", out, err)
	}
}

