package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"strings"
	"testing"
)

var tomlBytes = []byte(`names = ["John", "Frank"]

[person]
  name = "Tom"
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
		if !reflect.DeepEqual(tomlMap, got) {
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
	got, err := (&storage.TOMLParser{}).ToBytes(tomlMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if string(tomlBytes) != string(got) {
		t.Errorf("expected:\n%s\ngot:\n%s", tomlBytes, got)
	}
}
