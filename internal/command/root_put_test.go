package command_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"io/ioutil"
	"strings"
	"testing"
)

func putTest(in string, varType string, parser string, selector string, value string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"put", varType, "-p", parser, "-s", selector, value,
		}
		fmt.Println(args)

		cmd.SetOut(outputBuffer)
		cmd.SetIn(strings.NewReader(in))
		cmd.SetArgs(args)

		err := cmd.Execute()

		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}

		output, err := ioutil.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		out = strings.TrimSpace(out)
		outputStr := strings.TrimSpace(string(output))
		if out != outputStr {
			t.Errorf("expected result %v, got %v", out, outputStr)
		}
	}
}

func putStringTest(in string, parser string, selector string, value string, out string, expErr error) func(t *testing.T) {
	return putTest(in, "string", parser, selector, value, out, expErr)
}

func putIntTest(in string, parser string, selector string, value string, out string, expErr error) func(t *testing.T) {
	return putTest(in, "int", parser, selector, value, out, expErr)
}

func putBoolTest(in string, parser string, selector string, value string, out string, expErr error) func(t *testing.T) {
	return putTest(in, "bool", parser, selector, value, out, expErr)
}

func putStringTestForParserJSON(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`{
  "id": "x"
}`, parser, "id", "y", `{
  "id": "y"
}`, nil))
	}
}

func putIntTestForParserJSON(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`{
  "id": 123
}`, parser, "id", "456", `{
  "id": 456
}`, nil))
	}
}

func putBoolTestForParserJSON(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`{
  "id": true
}`, parser, "id", "false", `{
  "id": false
}`, nil))
	}
}

func putStringTestForParserYAML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`
id: "x"
`, parser, "id", "y", `
id: "y"
`, nil))
	}
}

func putIntTestForParserYAML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`
id: 123
`, parser, "id", "456", `
id: 456
`, nil))
	}
}

func putBoolTestForParserYAML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`
id: true
`, parser, "id", "false", `
id: false
`, nil))
	}
}

func putStringTestForParserTOML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`
id = "x"
`, parser, "id", "y", `
id = "y"
`, nil))
	}
}

func putIntTestForParserTOML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`
id = 123
`, parser, "id", "456", `
id = 456
`, nil))
	}
}

func putBoolTestForParserTOML(parser string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`
id = true
`, parser, "id", "false", `
id = false
`, nil))
	}
}