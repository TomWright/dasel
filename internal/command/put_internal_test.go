package command

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/internal/storage"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		Value     string
		ValueType string
		Out       interface{}
		Err       error
	}{
		{Value: "a", ValueType: "string", Out: "a"},
		{Value: "1", ValueType: "string", Out: "1"},
		{Value: "1", ValueType: "int", Out: int64(1)},
		{Value: "a", ValueType: "int", Err: fmt.Errorf("could not parse int [%s]: strconv.ParseInt: parsing \"%s\": invalid syntax", "a", "a")},
		{Value: "true", ValueType: "string", Out: "true"},
		{Value: "false", ValueType: "string", Out: "false"},
		{Value: "true", ValueType: "bool", Out: true},
		{Value: "false", ValueType: "bool", Out: false},
		{Value: "a", ValueType: "bool", Err: fmt.Errorf("could not parse bool [%s]: unhandled value", "a")},
		{Value: "a", ValueType: "bad", Err: fmt.Errorf("unhandled type: %s", "bad")},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(fmt.Sprintf("%s_%s", tc.Value, tc.ValueType), func(t *testing.T) {
			got, err := parseValue(tc.Value, tc.ValueType)
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if !reflect.DeepEqual(tc.Out, got) {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

func TestShouldReadFromStdin(t *testing.T) {
	if exp, got := false, shouldReadFromStdin("asd"); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := true, shouldReadFromStdin(""); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestGetReadParser(t *testing.T) {
	tests := []struct {
		File   string
		Parser string
		Out    storage.Parser
		Err    error
	}{
		{File: "a.json", Out: &storage.JSONParser{}},
		{Parser: "json", Out: &storage.JSONParser{}},
		{File: "a.yaml", Out: &storage.YAMLParser{}},
		{File: "a.yml", Out: &storage.YAMLParser{}},
		{Parser: "yaml", Out: &storage.YAMLParser{}},
		{File: "a.txt", Err: fmt.Errorf("could not get read parser from filename: unknown parser: .txt")},
		{Parser: "txt", Err: fmt.Errorf("could not get read parser: unknown parser: txt")},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run("Test", func(t *testing.T) {
			got, err := getReadParser(tc.File, tc.Parser, "")
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Out != got {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

func putTest(in string, parser string, selector string, value string, valueType string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		outputBuffer := bytes.NewBuffer([]byte{})

		err := runGenericPutCommand(genericPutOptions{
			Parser:    parser,
			Selector:  selector,
			Value:     value,
			ValueType: valueType,
			Reader:    strings.NewReader(in),
			Writer:    outputBuffer,
		}, nil)

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
		got := strings.TrimSpace(string(output))

		if out != got {
			t.Errorf("expected result:\n%s\ngot:\n%s", out, got)
		}
	}
}

func putObjectTest(in string, parser string, selector string, values []string, valueTypes []string, out string, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		outputBuffer := bytes.NewBuffer([]byte{})

		err := runPutObjectCommand(putObjectOpts{
			Parser:      parser,
			Selector:    selector,
			InputValues: values,
			InputTypes:  valueTypes,
			Reader:      strings.NewReader(in),
			Writer:      outputBuffer,
		}, nil)

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
		got := strings.TrimSpace(string(output))

		if out != got {
			t.Errorf("expected result:\n%s\ngot:\n%s", out, got)
		}
	}
}

func TestPut(t *testing.T) {
	t.Run("MissingParserFlag", func(t *testing.T) {
		err := runGenericPutCommand(genericPutOptions{}, nil)
		if err == nil || err.Error() != "read parser flag required when reading from stdin" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("InvalidVarType", func(t *testing.T) {
		err := runGenericPutCommand(genericPutOptions{Parser: "yaml", ValueType: "int", Value: "asd", Reader: bytes.NewBuffer([]byte{})}, nil)
		if err == nil || err.Error() != "could not parse int [asd]: strconv.ParseInt: parsing \"asd\": invalid syntax" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("FailedWrite", func(t *testing.T) {
		err := runGenericPutCommand(genericPutOptions{Parser: "yaml", ValueType: "string", Selector: ".name", Value: "asd", Reader: bytes.NewBuffer([]byte{}), Writer: &failingWriter{}}, nil)
		if err == nil || err.Error() != "could not write output: could not write to output file: could not write data: i am meant to fail at writing" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("ObjectMissingParserFlag", func(t *testing.T) {
		err := runPutObjectCommand(putObjectOpts{}, nil)
		if err == nil || err.Error() != "read parser flag required when reading from stdin" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("ObjectInvalidTypes", func(t *testing.T) {
		err := runPutObjectCommand(putObjectOpts{
			File:        "../../tests/assets/example.json",
			InputTypes:  []string{"string"},
			InputValues: []string{"x", "y"},
		}, nil)
		if err == nil || err.Error() != "exactly 2 types are required, got 1" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("JSON", func(t *testing.T) {
		t.Run("SingleProperty", putTest(`{
  "details": {
    "age": 27,
    "name": "Tom"
  },
  "id": "1111"
}`, "json", ".id", "2222", "string", `{
  "details": {
    "age": 27,
    "name": "Tom"
  },
  "id": "2222"
}`, nil))

		t.Run("ObjectPropertyString", putTest(`{
  "details": {
    "age": 27,
    "name": "Tom"
  },
  "id": "1111"
}`, "json", ".details.name", "Frank", "string", `{
  "details": {
    "age": 27,
    "name": "Frank"
  },
  "id": "1111"
}`, nil))

		t.Run("ObjectPropertyInt", putTest(`{
  "details": {
    "age": 20,
    "name": "Tom"
  },
  "id": "1111"
}`, "json", ".details.age", "27", "int", `{
  "details": {
    "age": 27,
    "name": "Tom"
  },
  "id": "1111"
}`, nil))

		t.Run("IndexString", putTest(`{
  "numbers": [
    "one",
    "two",
    "three"
  ]
}`, "json", ".numbers.[1]", "four", "string", `{
  "numbers": [
    "one",
    "four",
    "three"
  ]
}`, nil))

		t.Run("DynamicString", putTest(`{
  "numbers": [
	"one",
	"two",
	"three"
  ]
}`, "json", ".numbers.(value=three)", "four", "string", `{
  "numbers": [
    "one",
    "two",
    "four"
  ]
}`, nil))

		t.Run("DynamicInt", putTest(`{
  "numbers": [
	1,
	2,
	3
  ]
}`, "json", ".numbers.(value=3)", "4", "int", `{
  "numbers": [
    1,
    2,
    4
  ]
}`, nil))

		t.Run("IndexInt", putTest(`{
  "numbers": [
    1,
    2,
    3
  ]
}`, "json", ".numbers.[1]", "4", "int", `{
  "numbers": [
    1,
    4,
    3
  ]
}`, nil))

		t.Run("DynamicString", putTest(`{
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
    }
  ]
}`, "json", ".numbers.(number=two).number", "four", "string", `{
  "numbers": [
    {
      "number": "one",
      "rank": 1
    },
    {
      "number": "four",
      "rank": 2
    },
    {
      "number": "three",
      "rank": 3
    }
  ]
}`, nil))

		t.Run("DynamicInt", putTest(`{
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
}`, "json", ".numbers.(rank=2).rank", "4", "int", `{
  "numbers": [
    {
      "number": "one",
      "rank": 1
    },
    {
      "number": "two",
      "rank": 4
    },
    {
      "number": "three",
      "rank": 3
    }
  ]
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
}`, "json", ".numbers.[]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `{
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
      "number": "five",
      "rank": 5
    }
  ]
}`, nil))

	})

	t.Run("YAML", func(t *testing.T) {
		t.Run("SingleProperty", putTest(`
details:
  age: 27
  name: Tom
id: 1111
`, "yaml", ".id", "2222", "string", `
details:
  age: 27
  name: Tom
id: "2222"
`, nil))

		t.Run("ObjectPropertyString", putTest(`
details:
  age: 27
  name: Tom
id: 1111
`, "yaml", ".details.name", "Frank", "string", `details:
  age: 27
  name: Frank
id: 1111
`, nil))

		t.Run("ObjectPropertyInt", putTest(`
details:
  age: 20
  name: Tom
id: 1111
`, "yaml", ".details.age", "27", "int", `
details:
  age: 27
  name: Tom
id: 1111
`, nil))

		t.Run("IndexString", putTest(`
numbers:
- one
- two
- three
`, "yaml", ".numbers.[1]", "four", "string", `
numbers:
- one
- four
- three
`, nil))

		t.Run("DynamicString", putTest(`
numbers:
- one
- two
- three
`, "yaml", ".numbers.(value=three)", "four", "string", `
numbers:
- one
- two
- four
`, nil))

		t.Run("DynamicInt", putTest(`
numbers:
- 1
- 2
- 3
`, "yaml", ".numbers.(value=3)", "4", "int", `
numbers:
- 1
- 2
- 4
`, nil))

		t.Run("IndexInt", putTest(`
numbers:
- 1
- 2
- 3
`, "yaml", ".numbers.[1]", "4", "int", `
numbers:
- 1
- 4
- 3
`, nil))

		t.Run("DynamicString", putTest(`
numbers:
- number: one
  rank: 1
- number: two
  rank: 2
- number: three
  rank: 3
`, "yaml", ".numbers.(number=two).number", "four", "string", `
numbers:
- number: one
  rank: 1
- number: four
  rank: 2
- number: three
  rank: 3
`, nil))

		t.Run("DynamicInt", putTest(`
numbers:
- number: one
  rank: 1
- number: two
  rank: 2
- number: three
  rank: 3
`, "yaml", ".numbers.(rank=2).rank", "4", "int", `
numbers:
- number: one
  rank: 1
- number: two
  rank: 4
- number: three
  rank: 3
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

		t.Run("AppendObject", putObjectTest(`
numbers:
- number: one
  rank: 1
- number: two
  rank: 2
- number: three
  rank: 3
`, "yaml", ".numbers.[]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `
numbers:
- number: one
  rank: 1
- number: two
  rank: 2
- number: three
  rank: 3
- number: five
  rank: 5
`, nil))

	})

	t.Run("TOML", func(t *testing.T) {
		t.Run("SingleProperty", putTest(`
id = "1111"

[details]
  age = 27
  name = "Tom"
`, "toml", ".id", "2222", "string", `
id = "2222"

[details]
  age = 27
  name = "Tom"
`, nil))

		t.Run("ObjectPropertyString", putTest(`
id = "1111"

[details]
  age = 27
  name = "Tom"
`, "toml", ".details.name", "Frank", "string", `
id = "1111"

[details]
  age = 27
  name = "Frank"
`, nil))

		t.Run("ObjectPropertyInt", putTest(`
id = "1111"

[details]
  age = 20
  name = "Tom"
`, "toml", ".details.age", "27", "int", `
id = "1111"

[details]
  age = 27
  name = "Tom"
`, nil))

		t.Run("IndexString", putTest(`
numbers = ["one", "two", "three"]
`, "toml", ".numbers.[1]", "four", "string", `
numbers = ["one", "four", "three"]
`, nil))

		t.Run("DynamicString", putTest(`
numbers = ["one", "two", "three"]
`, "toml", ".numbers.(value=three)", "four", "string", `
numbers = ["one", "two", "four"]
`, nil))

		t.Run("DynamicInt", putTest(`
numbers = [1, 2, 3]
`, "toml", ".numbers.(value=3)", "4", "int", `
numbers = [1, 2, 4]
`, nil))

		t.Run("IndexInt", putTest(`
numbers = [1, 2, 3]
`, "toml", ".numbers.[1]", "4", "int", `
numbers = [1, 4, 3]
`, nil))

		t.Run("DynamicString", putTest(`
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, "toml", ".numbers.(number=two).number", "four", "string", `
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "four"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, nil))

		t.Run("DynamicInt", putTest(`
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, "toml", ".numbers.(rank=2).rank", "4", "int", `
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 4

[[numbers]]
  number = "three"
  rank = 3
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

		t.Run("AppendObject", putObjectTest(`
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3
`, "toml", ".numbers.[]", []string{"number=five", "rank=5"}, []string{"string", "int"}, `
[[numbers]]
  number = "one"
  rank = 1

[[numbers]]
  number = "two"
  rank = 2

[[numbers]]
  number = "three"
  rank = 3

[[numbers]]
  number = "five"
  rank = 5
`, nil))

	})
}
