package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"strings"
	"testing"
)

var yamlBytes = []byte(`name: Tom
`)
var yamlMap = map[string]interface{}{
	"name": "Tom",
}

var yamlBytesMulti = []byte(`name: Tom
---
name: Jim
`)
var yamlMapMulti = []interface{}{
	map[string]interface{}{
		"name": "Tom",
	},
	map[string]interface{}{
		"name": "Jim",
	},
}

func TestYAMLParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).FromBytes(yamlBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicSingleDocument{Value: yamlMap}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).FromBytes(yamlBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicMultiDocument{Values: yamlMapMulti}

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		_, err := (&storage.YAMLParser{}).FromBytes([]byte(`{1:asd`))
		if err == nil || !strings.Contains(err.Error(), "could not unmarshal data") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(nil, got) {
			t.Errorf("expected %v, got %v", nil, got)
		}
	})
}

func TestYAMLParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(yamlMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: yamlMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidMulti", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: yamlMapMulti})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytesMulti) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytesMulti, got)
		}
	})
}
