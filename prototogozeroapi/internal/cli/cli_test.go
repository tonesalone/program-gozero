package cli

import "testing"

func TestRun_NoArgs(t *testing.T) {
	err := Run([]string{})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

