package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestChangeDefaultCommand(t *testing.T) {
	cachedArgs := os.Args
	defer func() {
		os.Args = cachedArgs
	}()

	testArgs := func(in []string, exp []string) func(t *testing.T) {
		return func(t *testing.T) {
			os.Args = in

			cmd := command.NewRootCMD()
			command.ChangeDefaultCommand(cmd, "select")

			got := os.Args
			if !reflect.DeepEqual(exp, got) {
				t.Errorf("expected args %v, got %v", exp, got)
			}
		}
	}

	t.Run("ChangeToSelect", testArgs(
		[]string{"dasel", "-p", "json", ".name"},
		[]string{"dasel", "select", "-p", "json", ".name"},
	))

	t.Run("AlreadySelect", testArgs(
		[]string{"dasel", "select", "-p", "json", ".name"},
		[]string{"dasel", "select", "-p", "json", ".name"},
	))

	t.Run("AlreadyPut", testArgs(
		[]string{"dasel", "put", "-p", "json", "-t", "string", "name=Tom"},
		[]string{"dasel", "put", "-p", "json", "-t", "string", "name=Tom"},
	))
}

func TestRootCMD(t *testing.T) {
	expectErr := func(args []string, expErr string) func(t *testing.T) {
		return func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err == nil || !strings.Contains(err.Error(), expErr) {
				t.Errorf("unexpected error: %v", err)
				return
			}
		}
	}
	expectOutput := func(in string, args []string, exp string) func(t *testing.T) {
		return func(t *testing.T) {
			cmd := command.NewRootCMD()
			outputBuffer := bytes.NewBuffer([]byte{})

			cmd.SetIn(bytes.NewReader([]byte(in)))
			cmd.SetOut(outputBuffer)
			cmd.SetArgs(args)

			err := cmd.Execute()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := outputBuffer.String()
			if exp != got {
				t.Errorf("expected %s, got %s", exp, got)
			}
		}
	}

	t.Run("Select", func(t *testing.T) {
		t.Run("JSON", selectTestForParser("json", jsonData, jsonDataSingle))
		t.Run("YAML", selectTestForParser("yaml", yamlData, yamlDataSingle))
		t.Run("TOML", selectTestForParser("toml", tomlData, tomlDataSingle))
		t.Run("InvalidFile", expectErr(
			[]string{"select", "-f", "bad.json", "-s", "x"},
			"could not open input file",
		))
		t.Run("MissingParser", expectErr(
			[]string{"select", "-s", "x"},
			"parser flag required when reading from stdin",
		))
		t.Run("Stdin", expectOutput(
			`{"name": "Tom"}`,
			[]string{"select", "-f", "stdin", "-p", "json", "-s", ".name"},
			`Tom
`,
		))
		t.Run("StdinAlias", expectOutput(
			`{"name": "Tom"}`,
			[]string{"select", "-f", "-", "-p", "json", "-s", ".name"},
			`Tom
`,
		))
	})
	t.Run("PutString", func(t *testing.T) {
		t.Run("JSON", putStringTestForParserJSON())
		t.Run("YAML", putStringTestForParserYAML())
		t.Run("TOML", putStringTestForParserTOML())

		t.Run("InvalidFile", expectErr(
			[]string{"put", "string", "-f", "bad.json", "-s", "x", "y"},
			"could not open input file",
		))
		t.Run("MissingParser", expectErr(
			[]string{"put", "string", "-s", "x", "y"},
			"parser flag required when reading from stdin",
		))
		t.Run("StdinStdout", expectOutput(
			`{"name": "Tom"}`,
			[]string{"put", "string", "-f", "stdin", "-o", "stdout", "-p", "json", "-s", ".name", "Frank"},
			`{
  "name": "Frank"
}
`,
		))
		t.Run("StdinStdoutAlias", expectOutput(
			`{"name": "Tom"}`,
			[]string{"put", "string", "-f", "-", "-o", "-", "-p", "json", "-s", ".name", "Frank"},
			`{
  "name": "Frank"
}
`,
		))
	})
	t.Run("PutInt", func(t *testing.T) {
		t.Run("JSON", putIntTestForParserJSON())
		t.Run("YAML", putIntTestForParserYAML())
		t.Run("TOML", putIntTestForParserTOML())
	})
	t.Run("PutBool", func(t *testing.T) {
		t.Run("JSON", putBoolTestForParserJSON())
		t.Run("YAML", putBoolTestForParserYAML())
		t.Run("TOML", putBoolTestForParserTOML())
	})
	t.Run("PutObject", func(t *testing.T) {
		t.Run("JSON", putObjectTestForParserJSON())
		t.Run("YAML", putObjectTestForParserYAML())
		t.Run("TOML", putObjectTestForParserTOML())
	})
}
