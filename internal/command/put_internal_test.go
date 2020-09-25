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

func TestGetParser(t *testing.T) {
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
		{File: "a.txt", Err: fmt.Errorf("could not get parser from filename: unknown parser: .txt")},
		{Parser: "txt", Err: fmt.Errorf("could not get parser: unknown parser: txt")},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run("Test", func(t *testing.T) {
			got, err := getParser(tc.File, tc.Parser)
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
		})

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

		if out != string(output) {
			t.Errorf("expected result %v, got %v", out, string(output))
		}
	}
}

func TestPut(t *testing.T) {
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
    "age": 27,
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
	})
}
