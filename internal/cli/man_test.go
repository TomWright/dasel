package cli_test

import (
	"strings"
	"testing"
)

func TestMan(t *testing.T) {
	stdout, _, err := runDaselCmd([]string{"man"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stdout == "" {
		t.Fatal("expected non-empty output")
	}

	// Roff directives
	if !strings.Contains(stdout, ".TH") {
		t.Error("expected output to contain .TH roff directive")
	}
	if !strings.Contains(stdout, ".SH") {
		t.Error("expected output to contain .SH roff directive")
	}

	// Description
	if !strings.Contains(stdout, "Query and modify data structures") {
		t.Error("expected output to contain description text")
	}

	// Subcommands
	for _, sub := range []string{"query", "version", "completion", "man"} {
		if !strings.Contains(stdout, sub) {
			t.Errorf("expected output to contain subcommand %q", sub)
		}
	}

	// Key flags (roff uses \- for dashes)
	for _, flag := range []string{"in", "out", "compact"} {
		if !strings.Contains(stdout, flag) {
			t.Errorf("expected output to contain flag %q", flag)
		}
	}
}
