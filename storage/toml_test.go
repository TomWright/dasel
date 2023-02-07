package storage_test

import (
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
	"reflect"
	"strings"
	"testing"
)

var tomlBytes = []byte(`names = ['John', 'Frank']

[person]
name = 'Tom'
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
		if !reflect.DeepEqual(tomlMap, got.Interface()) {
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
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf(tomlMap))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(tomlBytes) != string(got) {
			t.Errorf("expected:\n---\n%s\n---\ngot:\n---\n%s\n---", tomlBytes, got)
		}
	})
	t.Run("SingleDocument", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf(tomlMap).WithMetadata("isSingleDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(tomlBytes) != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", tomlBytes, got)
		}
	})
	t.Run("SingleDocumentColourise", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf(tomlMap).WithMetadata("isSingleDocument", true), storage.ColouriseOption(true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		expBuf, _ := storage.Colourise(string(tomlBytes), "toml")
		exp := expBuf.Bytes()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", string(exp), string(got))
		}
	})
	t.Run("SingleDocumentCustomIndent", func(t *testing.T) {
		res, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf(tomlMap).WithMetadata("isSingleDocument", true), storage.IndentOption("   "))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(res)
		exp := `names = ['John', 'Frank']

[person]
   name = 'Tom'
`
		if exp != got {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("MultiDocument", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf([]interface{}{tomlMap, tomlMap}).WithMetadata("isMultiDocument", true))
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
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf("asd"))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `'asd'
`
		if exp != string(got) {
			t.Errorf("expected:\n---\n%s\n---\ngot:\n---\n%s\n---", exp, got)
		}
	})
	t.Run("DefaultValue", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf("asd"))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `'asd'
`
		if exp != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("MultiDocumentValue", func(t *testing.T) {
		got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf([]interface{}{"asd", 123}).WithMetadata("isMultiDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `'asd'
123
`
		if exp != string(got) {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	// t.Run("time.Time", func(t *testing.T) {
	// 	v, _ := time.Parse(time.RFC3339, "2022-01-02T12:34:56Z")
	// 	got, err := (&storage.TOMLParser{}).ToBytes(dasel.ValueOf(v))
	// 	if err != nil {
	// 		t.Errorf("unexpected error: %s", err)
	// 		return
	// 	}
	// 	exp := `2022-01-02T12:34:56Z
	// `
	// 		if exp != string(got) {
	// 			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	// 		}
	// 	})
}
