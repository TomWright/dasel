package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRootCMD_Put(t *testing.T) {
	t.Run("String", func(t *testing.T) {
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
	t.Run("Int", func(t *testing.T) {
		t.Run("JSON", putIntTestForParserJSON())
		t.Run("YAML", putIntTestForParserYAML())
		t.Run("TOML", putIntTestForParserTOML())
	})
	t.Run("Bool", func(t *testing.T) {
		t.Run("JSON", putBoolTestForParserJSON())
		t.Run("YAML", putBoolTestForParserYAML())
		t.Run("TOML", putBoolTestForParserTOML())
	})
	t.Run("Object", func(t *testing.T) {
		t.Run("JSON", putObjectTestForParserJSON())
		t.Run("YAML", putObjectTestForParserYAML())
		t.Run("TOML", putObjectTestForParserTOML())
	})
}

func putTest(in string, varType string, parser string, selector string, value string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"put", varType,
		}
		args = append(args, additionalArgs...)
		args = append(args, "-p", parser, selector, value)

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

func putObjectTest(in string, parser string, selector string, values []string, types []string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"put", "object", "-p", parser, "-o", "stdout", selector,
		}
		for _, t := range types {
			args = append(args, "-t", t)
		}
		args = append(args, values...)

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

func putStringTestForParserJSON() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`{
  "id": "x"
}`, "json", "id", "y", `{
  "id": "y"
}`, nil))
	}
}

func putIntTestForParserJSON() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`{
  "id": 123
}`, "json", "id", "456", `{
  "id": 456
}`, nil))
	}
}

func putBoolTestForParserJSON() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`{
  "id": true
}`, "json", "id", "false", `{
  "id": false
}`, nil))
	}
}

func putObjectTestForParserJSON() func(t *testing.T) {
	return func(t *testing.T) {

		t.Run("OverwriteObject", putObjectTest(`{
  "numbers": [
    {
      "rank": 1,
      "number": "one"
    },
    {
      "rank": 2,
      "number": "two"
    },
    {
      "rank": 3,
      "number": "three"
    }
  ]
}`, "json", ".numbers.[0]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `{
  "numbers": [
    {
      "number": "five",
      "rank": 5
    },
    {
      "number": "two",
      "rank": 2
    },
    {
      "number": "three",
      "rank": 3
    }
  ]
}`, nil))

	}
}

func putStringTestForParserYAML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`
id: "x"
`, "yaml", "id", "y", `
id: "y"
`, nil))
	}
}

func putIntTestForParserYAML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`
id: 123
`, "yaml", "id", "456", `
id: 456
`, nil))
	}
}

func putBoolTestForParserYAML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`
id: true
`, "yaml", "id", "false", `
id: false
`, nil))
	}
}

func putObjectTestForParserYAML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("OverwriteObject", putObjectTest(`
numbers:
- number: one
  rank: 1
- number: two
  rank: 2
- number: three
  rank: 3
`, "yaml", ".numbers.[0]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `
numbers:
- number: five
  rank: 5
- number: two
  rank: 2
- number: three
  rank: 3
`, nil))
	}
}

func putStringTestForParserTOML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putStringTest(`
id = "x"
`, "toml", "id", "y", `
id = "y"
`, nil))
	}
}

func putIntTestForParserTOML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putIntTest(`
id = 123
`, "toml", "id", "456", `
id = 456
`, nil))
	}
}

func putBoolTestForParserTOML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("Property", putBoolTest(`
id = true
`, "toml", "id", "false", `
id = false
`, nil))
	}
}

func putObjectTestForParserTOML() func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("OverwriteObject", putObjectTest(`
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, "toml", ".numbers.[0]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `
[[numbers]]
  number = "five"
  rank = 5

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, nil))
	}
}
