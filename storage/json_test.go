package storage_test

import (
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/dencoding"
	"github.com/tomwright/dasel/v2/storage"
	"reflect"
	"testing"
)

var jsonBytes = []byte(`{
  "name": "Tom"
}
`)
var jsonMap = dencoding.NewMap().Set("name", "Tom")

func TestJSONParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := jsonMap
		if !reflect.DeepEqual(exp, got.Interface()) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := jsonMapMulti
		if !reflect.DeepEqual(exp, got.Interface()) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("ValidMultiDocumentMixed", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes(jsonBytesMultiMixed)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := jsonMapMultiMixed
		if !reflect.DeepEqual(exp, got.Interface()) {
			t.Errorf("expected %v, got %v", jsonMap, got)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(dasel.Value{}, got) {
			t.Errorf("expected %v, got %v", nil, got)
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
		got, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMap))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMap).WithMetadata("isSingleDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
		}
	})

	t.Run("ValidSingleNoPrettyPrint", func(t *testing.T) {
		res, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMap).WithMetadata("isSingleDocument", true), storage.PrettyPrintOption(false))
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

	t.Run("ValidSingleColourise", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMap).WithMetadata("isSingleDocument", true), storage.ColouriseOption(true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		expBuf, _ := storage.Colourise(`{
  "name": "Tom"
}
`, "json")
		exp := expBuf.Bytes()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})

	t.Run("ValidSingleCustomIndent", func(t *testing.T) {
		res, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMap).WithMetadata("isSingleDocument", true), storage.IndentOption("   "))
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
		got, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMapMulti).WithMetadata("isMultiDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(jsonBytesMulti) != string(got) {
			t.Errorf("expected %v, got %v", string(jsonBytesMulti), string(got))
		}
	})

	t.Run("ValidMultiMixed", func(t *testing.T) {
		got, err := (&storage.JSONParser{}).ToBytes(dasel.ValueOf(jsonMapMultiMixed).WithMetadata("isMultiDocument", true))
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

var jsonMapMulti = []any{
	dencoding.NewMap().Set("name", "Tom"),
	dencoding.NewMap().Set("name", "Ellis"),
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
	dencoding.NewMap().Set("name", "Tom").Set("other", true),
	dencoding.NewMap().Set("name", "Ellis"),
}
