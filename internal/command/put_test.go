package command

import (
	"fmt"
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
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
