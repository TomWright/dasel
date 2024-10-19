package cli_test

import (
	"os"
	"testing"
)

func TestManCommand(t *testing.T) {
	t.Skip("Temporarily disabled")
	tempDir := t.TempDir()

	_, _, err := runDasel([]string{"man", "-o", tempDir}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedFiles := []string{
		"dasel-completion-bash.1",
		"dasel-completion-fish.1",
		"dasel-completion-powershell.1",
		"dasel-completion-zsh.1",
		"dasel-completion.1",
		//"dasel-delete.1",
		"dasel-man.1",
		//"dasel-put.1",
		//"dasel-validate.1",
		"dasel.1",
	}

	if len(files) != len(expectedFiles) {
		t.Fatalf("expected %d files, got %d", len(expectedFiles), len(files))
	}

	for i, f := range files {
		if f.Name() != expectedFiles[i] {
			t.Fatalf("expected %v, got %v", expectedFiles[i], f.Name())
		}
	}
}
