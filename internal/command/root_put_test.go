package command_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/command"
	"io/ioutil"
	"strings"
	"testing"
)

func TestRootCMD_Put(t *testing.T) {
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

func TestRootCMD_Put_JSON(t *testing.T) {
	t.Run("String", putStringTest(`{
  "id": "x"
}`, "json", "id", "y", `{
  "id": "y"
}`, nil))
	t.Run("Int", putIntTest(`{
  "id": 123
}`, "json", "id", "456", `{
  "id": 456
}`, nil))
	t.Run("Bool", putBoolTest(`{
  "id": true
}`, "json", "id", "false", `{
  "id": false
}`, nil))
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

func TestRootCMD_Put_YAML(t *testing.T) {
	t.Run("String", putStringTest(`
id: "x"
name: "Tom"
`, "yaml", "id", "y", `
id: "y"
name: Tom
`, nil))
	t.Run("Int", putIntTest(`
id: 123
`, "yaml", "id", "456", `
id: 456
`, nil))
	t.Run("Bool", putBoolTest(`
id: true
`, "yaml", "id", "false", `
id: false
`, nil))
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
	t.Run("StringInMultiDocument", putStringTest(`
id: "x"
---
id: "y"
---
id: "z"
`, "yaml", "[1].id", "1", `
id: x
---
id: "1"
---
id: z
`, nil))
}

func TestRootCMD_Put_TOML(t *testing.T) {
	t.Run("String", putStringTest(`
id = "x"
`, "toml", "id", "y", `
id = "y"
`, nil))
	t.Run("Int", putIntTest(`
id = 123
`, "toml", "id", "456", `
id = 456
`, nil))
	t.Run("Bool", putBoolTest(`
id = true
`, "toml", "id", "false", `
id = false
`, nil))
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

func TestRootCMD_Put_XML(t *testing.T) {
	t.Run("String", putStringTest(`<data><id>x</id></data>`, "xml", ".data.id", "y", `<data>
  <id>y</id>
</data>
`, nil))
	t.Run("Int", putIntTest(`<data><id>1</id></data>`, "xml", ".data.id", "2", `<data>
  <id>2</id>
</data>
`, nil))
	t.Run("Bool", putBoolTest(`<data><id>false</id></data>`, "xml", ".data.id", "true", `<data>
  <id>true</id>
</data>
`, nil))
	t.Run("OverwriteObject", putObjectTest(`<data><id>x</id></data>`, "xml", ".data", []string{"id=y", "rank=5"}, []string{"string", "int"}, `<data>
  <id>y</id>
  <rank>5</rank>
</data>
`, nil))
	t.Run("AppendObject", putObjectTest(`<data><item><value>1</value></item><item><value>2</value></item></data>`, "xml", ".data.item.[]", []string{"value=3"}, []string{"int"}, `<data>
  <item>
    <value>1</value>
  </item>
  <item>
    <value>2</value>
  </item>
  <item>
    <value>3</value>
  </item>
</data>
`, nil))
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
