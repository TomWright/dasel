package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/tomwright/dasel/v3/internal/cli"
)

func runDaselCmd(args []string) (string, string, error) {
	stdOut := bytes.NewBuffer(nil)
	stdErr := bytes.NewBuffer(nil)
	stdIn := bytes.NewReader(nil)

	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = append([]string{"dasel"}, args...)
	_, err := cli.Run(stdIn, stdOut, stdErr)
	return stdOut.String(), stdErr.String(), err
}

func TestCompletion(t *testing.T) {
	tests := []struct {
		shell   string
		markers []string
	}{
		{
			shell: "bash",
			markers: []string{
				"_dasel()",
				"complete -F _dasel dasel",
				"COMPREPLY",
				"compgen",
			},
		},
		{
			shell: "zsh",
			markers: []string{
				"#compdef dasel",
				"_dasel()",
				"_arguments",
				"_describe",
			},
		},
		{
			shell: "fish",
			markers: []string{
				"complete -c dasel",
				"__fish_use_subcommand",
				"__fish_seen_subcommand_from",
			},
		},
		{
			shell: "powershell",
			markers: []string{
				"Register-ArgumentCompleter",
				"-CommandName dasel",
				"CompletionResult",
			},
		},
	}

	subcommands := []string{"query", "version", "completion", "man"}
	// Fish uses -l flag_name format, others use --flag_name
	flags := []string{"in", "out", "compact"}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			stdout, _, err := runDaselCmd([]string{"completion", tt.shell})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if stdout == "" {
				t.Fatal("expected non-empty output")
			}
			if !strings.Contains(stdout, "dasel") {
				t.Error("expected output to contain 'dasel'")
			}
			for _, marker := range tt.markers {
				if !strings.Contains(stdout, marker) {
					t.Errorf("expected output to contain %q", marker)
				}
			}
			for _, sub := range subcommands {
				if !strings.Contains(stdout, sub) {
					t.Errorf("expected output to contain subcommand %q", sub)
				}
			}
			for _, flag := range flags {
				if !strings.Contains(stdout, flag) {
					t.Errorf("expected output to contain flag %q", flag)
				}
			}
		})
	}
}
