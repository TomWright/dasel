package command_test

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"io"
	"os"
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

	t.Run("InvalidSingleSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"put", "string", "-f", "stdin", "-o", "stdout", "-p", "json", "-s", "[-]", "Frank"},
		"invalid index: -",
	))
	t.Run("InvalidMultiSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"put", "string", "-f", "stdin", "-o", "stdout", "-p", "json", "-m", "-s", "[-]", "Frank"},
		"invalid index: -",
	))

	t.Run("InvalidObjectSingleSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"put", "object", "-f", "stdin", "-o", "stdout", "-p", "json", "-t", "string", "-s", "[-]", "Frank"},
		"invalid index: -",
	))
	t.Run("InvalidMultiSelector", expectErrFromInput(
		`{"name": "Tom"}`,
		[]string{"put", "object", "-f", "stdin", "-o", "stdout", "-p", "json", "-m", "-t", "string", "-s", "[-]", "Frank"},
		"invalid index: -",
	))
}

func putTest(in string, varType string, parser string, selector string, value string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	args := []string{
		"put", varType,
	}
	args = append(args, additionalArgs...)
	args = append(args, "-p", parser, selector, value)

	return baseTest(in, out, expErr, args...)
}

func baseTest(in string, out string, expErr error, args ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

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

		output, err := io.ReadAll(outputBuffer)
		if err != nil {
			t.Errorf("unexpected error reading output buffer: %s", err)
			return
		}

		outputStr := string(output)
		if out != outputStr {
			t.Errorf("expected result %v, got %v", out, outputStr)
		}
	}
}

func putFileTest(in string, varType string, parser string, selector string, value string, out string, outFile string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		defer func() {
			_ = os.Remove(outFile)
		}()
		cmd := command.NewRootCMD()

		args := []string{
			"put", varType,
		}
		args = append(args, additionalArgs...)
		args = append(args, "-p", parser, "-o", outFile, selector, value)

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

		output, err := os.ReadFile(outFile)
		if err != nil {
			t.Errorf("could not read output file: %s", err)
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
}
`, nil))
	t.Run("Int", putIntTest(`{
  "id": 123
}`, "json", "id", "456", `{
  "id": 456
}
`, nil))
	t.Run("Bool", putBoolTest(`{
  "id": true
}`, "json", "id", "false", `{
  "id": false
}
`, nil))
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

	t.Run("OverwriteObjectAtRoot", putObjectTest(`{
  "rank": 1,
  "number": "one"
}`, "json", ".", []string{"number=five", "rank=5"}, []string{"string", "int"}, `{
  "number": "five",
  "rank": 5
}`, nil))

	t.Run("OverwriteObjectAtRootCompact", putObjectTest(`{
  "rank": 1,
  "number": "one"
}`, "json", ".", []string{"number=five", "rank=5"}, []string{"string", "int"}, `{"number":"five","rank":5}
`, nil, "-c"))

	t.Run("MultipleObject", putObjectTest(`{
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
}`, "json", ".numbers.[*]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `{
  "numbers": [
    {
      "number": "five",
      "rank": 5
    },
    {
      "number": "five",
      "rank": 5
    },
    {
      "number": "five",
      "rank": 5
    }
  ]
}`, nil, "-m"))

	t.Run("AppendObject", putObjectTest(`{
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
}`, "json", ".numbers.[]", []string{"rank=4", "number=four"}, []string{"int", "string"}, `{
  "numbers": [
    {
      "number": "one",
      "rank": 1
    },
    {
      "number": "two",
      "rank": 2
    },
    {
      "number": "three",
      "rank": 3
    },
    {
      "number": "four",
      "rank": 4
    }
  ]
}`, nil))

	t.Run("EmptyObject", putObjectTest(`{
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
}`, "json", ".numbers.[]", []string{}, []string{}, `{
  "numbers": [
    {
      "number": "one",
      "rank": 1
    },
    {
      "number": "two",
      "rank": 2
    },
    {
      "number": "three",
      "rank": 3
    },
    {}
  ]
}`, nil))

	t.Run("AppendObjectMulti", putObjectTest(`{
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
}`, "json", ".numbers.[]", []string{"rank=4", "number=four"}, []string{"int", "string"}, `{
  "numbers": [
    {
      "number": "one",
      "rank": 1
    },
    {
      "number": "two",
      "rank": 2
    },
    {
      "number": "three",
      "rank": 3
    },
    {
      "number": "four",
      "rank": 4
    }
  ]
}`, nil, "-m"))

	t.Run("MultipleString", putStringTest(`[
  {"value": "A"},
  {"value": "B"},
  {"value": "C"}
]`, "json", "[*].value", "X", `[
  {
    "value": "X"
  },
  {
    "value": "X"
  },
  {
    "value": "X"
  }
]
`, nil, "-m"))

	t.Run("KeySearch", putStringTest(`{
  "users": [
	{
	  "primary": true,
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  }
	},
	{
	  "primary": false,
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
	  "name": {
		"first": "Jim",
		"last": "Wright"
	  }
	}
  ]
}`, "json", ".(?:-=name).first", "Bobby", `{
  "users": [
    {
      "name": {
        "first": "Bobby",
        "last": "Wright"
      },
      "primary": true
    },
    {
      "extra": {
        "name": {
          "first": "Bobby",
          "last": "Blogs"
        }
      },
      "name": {
        "first": "Bobby",
        "last": "Wright"
      },
      "primary": false
    }
  ]
}
`, nil, "-m"))

	t.Run("ValueSearch", putStringTest(`{
  "users": [
	{
	  "primary": true,
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  }
	},
	{
	  "primary": false,
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
	  "name": {
		"first": "Jim",
		"last": "Wright"
	  }
	}
  ]
}`, "json", ".(?:.=Wright)", "Wrighto", `{
  "users": [
    {
      "name": {
        "first": "Tom",
        "last": "Wrighto"
      },
      "primary": true
    },
    {
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
      "name": {
        "first": "Jim",
        "last": "Wrighto"
      },
      "primary": false
    }
  ]
}
`, nil, "-m"))

	t.Run("KeyValueSearch", putStringTest(`{
  "users": [
	{
	  "primary": true,
	  "name": {
		"first": "Tom",
		"last": "Wright"
	  }
	},
	{
	  "primary": false,
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
	  "name": {
		"first": "Jim",
		"last": "Wright"
	  }
	}
  ]
}`, "json", ".(?:.last=Wright).first", "Fred", `{
  "users": [
    {
      "name": {
        "first": "Fred",
        "last": "Wright"
      },
      "primary": true
    },
    {
      "extra": {
        "name": {
          "first": "Joe",
          "last": "Blogs"
        }
      },
      "name": {
        "first": "Fred",
        "last": "Wright"
      },
      "primary": false
    }
  ]
}
`, nil, "-m"))

	t.Run("StringMultiObjectDocument", putStringTest(`{
  "id": "x"
}
{
  "id": "y"
}`, "json", ".[0].id", "z", `{
  "id": "z"
}
{
  "id": "y"
}
`, nil))

	t.Run("StringMultiArrayDocument", putStringTest(`[
  "a",
  "b",
  "c"
]
[
  "d",
  "e",
  "f"
]`, "json", ".[1].[1]", "z", `[
  "a",
  "b",
  "c"
]
[
  "d",
  "z",
  "f"
]
`, nil))

	t.Run("InsertDocumentAtProperty", putDocumentTest(`{}`, "json", ".person", `{"name":"Tom"}`, `{
  "person": {
    "name": "Tom"
  }
}
`, nil))

	t.Run("InsertDocumentAtPropertyCompact", putDocumentTest(`{}`, "json", ".person", `{"name":"Tom"}`, `{"person":{"name":"Tom"}}
`, nil, "-c"))

	t.Run("InvalidDocumentParser", putDocumentTest(`{}`, "json", ".person", `name: Tom`, ``,
		fmt.Errorf("could not get document parser: unknown parser: bad"), "-d", "bad"))

	t.Run("InsertDocumentAtPropertyWithDifferentParser", putDocumentTest(`{}`, "json", ".person", `name: Tom`, `{
  "person": {
    "name": "Tom"
  }
}
`, nil, "-d", "yaml"))

	t.Run("AppendDocumentToProperty", putDocumentTest(`{"people": []}`, "json", ".people.[]", `{"name":"Tom"}`, `{
  "people": [
    {
      "name": "Tom"
    }
  ]
}
`, nil))

	t.Run("AppendDocumentToPropertyMulti", putDocumentTest(`{"people": []}{"people": []}`, "json", ".[*].people.[]", `{"name":"Tom"}`, `{
  "people": [
    {
      "name": "Tom"
    }
  ]
}
{
  "people": [
    {
      "name": "Tom"
    }
  ]
}
`, nil, "-m"))

	t.Run("InsertDocumentAtRoot", putDocumentTest(`{}`, "json", ".", `{"name":"Tom"}`, `{
  "name": "Tom"
}
`, nil))

	t.Run("AppendDocumentAtRoot", putDocumentTest(`[]`, "json", ".[]", `{"name":"Tom"}`, `[
  {
    "name": "Tom"
  }
]
`, nil))

	t.Run("AppendDocumentAtRootMulti", putDocumentTest(`[][]`, "json", ".[*].[]", `{"name":"Tom"}`, `[
  {
    "name": "Tom"
  }
]
[
  {
    "name": "Tom"
  }
]
`, nil, "-m"))

	// https://github.com/TomWright/dasel/issues/66
	t.Run("PutJSONDocumentAtYAMLProperty", putDocumentTest(`foo: true
bar: 5
baz:
  qux: false
  quux: "yes"
  quuz: 7`, "yaml", ".baz", `{"qux": false,"quux": "no","quuz": 8}`, `bar: 5
baz:
  quux: "no"
  quuz: 8
  qux: false
foo: true
`, nil, "-d", "json"))

	t.Run("MergeInputDocumentsPut", putIntTest(`
{
"number": 1
}
{
"number": 2
}
{
"number": 3
}
`, "json", ".[0].number", `4`, `[
  {
    "number": 4
  },
  {
    "number": 2
  },
  {
    "number": 3
  }
]
`, nil, "--merge-input-documents"))

	t.Run("ValueFlag", func(t *testing.T) {
		// Test -v/--value flag
		// Workaround for https://github.com/TomWright/dasel/issues/117

		t.Run("StringWithDash", baseTest(`{
  "id": "x"
}`, `{
  "id": "-abc"
}
`, nil, "put", "string", "-p", "json", "-v", "-abc", ".id"))

		t.Run("NegativeInt", baseTest(`{
  "id": 1
}`, `{
  "id": -1
}
`, nil, "put", "int", "-p", "json", "-v", "-1", ".id"))
	})
	t.Run("StringWithDashWithSelectorFlag", baseTest(`{
  "id": "x"
}`, `{
  "id": "-abc"
}
`, nil, "put", "string", "-p", "json", "-v", "-abc", "-s", ".id"))

	t.Run("NegativeIntWithSelectorFlag", baseTest(`{
  "id": 1
}`, `{
  "id": -1
}
`, nil, "put", "int", "-p", "json", "-v", "-1", "-s", ".id"))
}

func TestRootCMD_Put_YAML(t *testing.T) {
	t.Run("String", putStringTest(`
id: "x"
name: "Tom"
`, "yaml", "id", "y", `id: "y"
name: Tom
`, nil))
	t.Run("StringInFile", putFileTest(`
id: "x"
name: "Tom"
`, "string", "yaml", "id", "y", `id: "y"
name: Tom
`, "TestRootCMD_Put_YAML_out.yaml", nil))
	t.Run("Int", putIntTest(`
id: 123
`, "yaml", "id", "456", `id: 456
`, nil))
	t.Run("Bool", putBoolTest(`
id: true
`, "yaml", "id", "false", `id: false
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
`, "yaml", "[1].id", "1", `id: x
---
id: "1"
---
id: z
`, nil))

	t.Run("StringWithDotInName", putStringTest(`
id: "asd"
my.name: "Tom"
`, "yaml", `my\.name`, "Jim", `id: asd
my.name: Jim
`, nil))

	t.Run("NewValInExistingMap", putStringTest(`
metadata:
  annotations:
    node.longhorn.io/default-disks-config: '[ { "name":"fast",  "path":"/mnt/data-fast1", "allowScheduling":true, "tags":["fast"]}, { "name":"slow",  "path":"/mnt/data-slow1", "allowScheduling":true, "tags":["slow"]} ]'
`, "yaml", `metadata.labels.node\.longhorn\.io\/create-default-disk`, "config", `metadata:
  annotations:
    node.longhorn.io/default-disks-config: '[ { "name":"fast",  "path":"/mnt/data-fast1",
      "allowScheduling":true, "tags":["fast"]}, { "name":"slow",  "path":"/mnt/data-slow1",
      "allowScheduling":true, "tags":["slow"]} ]'
  labels:
    node.longhorn.io/create-default-disk: config
`, nil))

	// https://github.com/TomWright/dasel/issues/102
	// Worked in v1.13.2
	t.Run("BlankInput", putStringTest(``, "yaml", `[0].job_name`, "logging", `- job_name: logging
`, nil))
}

func TestRootCMD_Put_TOML(t *testing.T) {
	t.Run("String", putStringTest(`
id = "x"
`, "toml", "id", "y", `id = "y"
`, nil))
	t.Run("Int", putIntTest(`
id = 123
`, "toml", "id", "456", `id = 456
`, nil))
	t.Run("Bool", putBoolTest(`
id = true
`, "toml", "id", "false", `id = false
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

func TestRootCMD_Put_CSV(t *testing.T) {
	t.Run("String", putStringTest(`id,name
1,Tom
2,Jim
`, "csv", ".[0].id", "3", `id,name
3,Tom
2,Jim
`, nil))
	t.Run("NewString", putStringTest(`id,name
1,Tom
2,Jim
`, "csv", ".[0].age", "27", `id,name,age
1,Tom,27
2,Jim,
`, nil))
}

func putObjectTest(in string, parser string, selector string, values []string, types []string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"put", "object", "-p", parser, "-o", "stdout",
		}
		for _, t := range types {
			args = append(args, "-t", t)
		}
		args = append(args, additionalArgs...)
		args = append(args, selector)
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

		output, err := io.ReadAll(outputBuffer)
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

func putDocumentTest(in string, parser string, selector string, document string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := command.NewRootCMD()
		outputBuffer := bytes.NewBuffer([]byte{})

		args := []string{
			"put", "document", "-p", parser, "-o", "stdout",
		}
		args = append(args, additionalArgs...)
		args = append(args, selector)
		args = append(args, document)

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
		if expErr != nil {
			return
		}

		output, err := io.ReadAll(outputBuffer)
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

func putStringTest(in string, parser string, selector string, value string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return putTest(in, "string", parser, selector, value, out, expErr, additionalArgs...)
}

func putIntTest(in string, parser string, selector string, value string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return putTest(in, "int", parser, selector, value, out, expErr, additionalArgs...)
}

func putBoolTest(in string, parser string, selector string, value string, out string, expErr error, additionalArgs ...string) func(t *testing.T) {
	return putTest(in, "bool", parser, selector, value, out, expErr, additionalArgs...)
}
