package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var jsonBytes = []byte(`{
  "name": "Tom"
}
`)
var jsonMap = map[string]interface{}{
	"name": "Tom",
}

func TestJSONParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicSingleDocument{Value: jsonMap}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicMultiDocument{
			Values: jsonMapMulti,
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("ValidMultiDocumentMixed", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytesMultiMixed)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicMultiDocument{
			Values: jsonMapMultiMixed,
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicSingleDocument{
			Value: map[string]interface{}{},
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
}

func TestJSONParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.JSONParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestJSONParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(jsonMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingleNoPrettyPrint", func(t *testing.T) {
		res, err := (&storage.JSONParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap}, storage.PrettyPrintOption(false))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(res)
		exp := `{"name":"Tom"}
`
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidSingleCustomIndent", func(t *testing.T) {
		res, err := (&storage.JSONParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap}, storage.IndentOption("   "))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		got := string(res)
		exp := `{
   "name": "Tom"
}
`
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidMulti", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(&storage.BasicMultiDocument{Values: jsonMapMulti})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytesMulti) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytesMulti), string(got))
		}
	})

	t.Run("ValidMultiMixed", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(&storage.BasicMultiDocument{Values: jsonMapMultiMixed})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytesMultiMixed) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytesMultiMixed), string(got))
		}
	})
}

var jsonBytesMulti = []byte(`{
  "name": "Tom"
}
{
  "name": "Ellis"
}
`)

var jsonMapMulti = []interface{}{
	map[string]interface{}{"name": "Tom"},
	map[string]interface{}{"name": "Ellis"},
}

var jsonBytesMultiMixed = []byte(`{
  "name": "Tom",
  "other": true
}
{
  "name": "Ellis"
}
`)

var jsonMapMultiMixed = []interface{}{
	map[string]interface{}{"name": "Tom", "other": true},
	map[string]interface{}{"name": "Ellis"},
}
