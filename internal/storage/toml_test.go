package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"strings"
	"testing"
)

var tomlBytes = []byte(`"names" = ["John", "Frank"]

["person"]
  "name" = "Tom"
`)
var tomlMap = map[string]interface{}{
	"person": map[string]interface{}{
		"name": "Tom",
	},
	"names": []interface{}{"John", "Frank"},
}

func TestTOMLParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).FromBytes(tomlBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(&storage.BasicSingleDocument{Value: tomlMap}, got) {
			t.Errorf("expected %v, got %v", tomlMap, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		_, err := (&storage.TOMLParser{}).FromBytes([]byte(`x:x`))
		if err == nil || !strings.Contains(err.Error(), "could not unmarshal data") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

func TestTOMLParser_ToBytes(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(tomlMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(tomlBytes) != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", tomlBytes, got)
		}
	})
	t.Run("SingleDocument", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: tomlMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(tomlBytes) != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", tomlBytes, got)
		}
	})
	t.Run("MultiDocument", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{tomlMap, tomlMap}})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := append([]byte{}, tomlBytes...)
		exp = append(exp, tomlBytes...)
		if string(exp) != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("SingleDocumentValue", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: "asd"})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
`
		if exp != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("DefaultValue", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes("asd")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
`
		if exp != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("MultiDocumentValue", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{"asd", "123"}})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
123
`
		if exp != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
}
