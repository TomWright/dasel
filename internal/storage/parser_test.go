package storage_test

import (
	"bytes"
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"strings"
	"testing"
)

func TestUnknownParserErr_Error(t *testing.T) {
	if exp, got := "unknown parser: x", (&storage.UnknownParserErr{Parser: "x"}).Error(); exp != got {
		t.Errorf("expected error %s, got %s", exp, got)
	}
}

func TestNewParserFromString(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "json", Out: &storage.JSONParser{}},
		{In: "yaml", Out: &storage.YAMLParser{}},
		{In: "toml", Out: &storage.TOMLParser{}},
		{In: "bad", Out: nil, Err: &storage.UnknownParserErr{Parser: "bad"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewParserFromString(tc.In)
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

func TestNewParserFromFilename(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "a.json", Out: &storage.JSONParser{}},
		{In: "a.yaml", Out: &storage.YAMLParser{}},
		{In: "a.yml", Out: &storage.YAMLParser{}},
		{In: "a.toml", Out: &storage.TOMLParser{}},
		{In: "a.txt", Out: nil, Err: &storage.UnknownParserErr{Parser: ".txt"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewParserFromFilename(tc.In)
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

var jsonData = map[string]interface{}{
	"name": "Tom",
	"preferences": map[string]interface{}{
		"favouriteColour": "red",
	},
	"colours": []interface{}{"red", "green", "blue"},
	"colourCodes": []interface{}{
		map[string]interface{}{
			"name": "red",
			"rgb": "ff0000",
		},
		map[string]interface{}{
			"name": "green",
			"rgb": "00ff00",
		},
		map[string]interface{}{
			"name": "blue",
			"rgb": "0000ff",
		},
	},
}

func TestLoadFromFile(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		data, err := storage.LoadFromFile("../../tests/assets/example.json", &storage.JSONParser{})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(jsonData, data) {
			t.Errorf("data does not match")
		}
	})
	t.Run("BaseFilePath", func(t *testing.T) {
		_, err := storage.LoadFromFile("x.json", &storage.JSONParser{})
		if err == nil || !strings.Contains(err.Error(), "could not open file") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

func TestWrite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		if err := storage.Write(&storage.JSONParser{}, map[string]interface{}{"name": "Tom"}, &buf); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if exp, got := `{
  "name": "Tom"
}`, buf.String(); exp != got {
			t.Errorf("unexpected output:\n%s\ngot:\n%s", exp, got)
		}
	})
}
