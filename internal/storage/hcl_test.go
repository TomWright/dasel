package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var hclBytes = []byte(`
baz = false
foo = "bar"
x = 500
`)
var hclMap = map[string]interface{}{
	"baz": false,
	"foo": "bar",
	"x":   float64(500),
}

func TestHCLParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.HCLParser{}).FromBytes(hclBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicSingleDocument{Value: hclMap}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&storage.HCLParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(nil, got) {
			t.Errorf("expected %v, got %v", nil, got)
		}
	})
}

func TestHCLParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.HCLParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestHCLParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.HCLParser{}).ToBytes(jsonMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(hclBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(hclBytes), string(got))
		}
	})

	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&storage.HCLParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(hclBytes) != string(got) {
			t.Errorf("expected %v, got %v", string(hclBytes), string(got))
		}
	})

	t.Run("ValidSingleNoPrettyPrint", func(t *testing.T) {
		res, err := (&storage.HCLParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap}, storage.PrettyPrintOption(false))
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
		res, err := (&storage.HCLParser{}).ToBytes(&storage.BasicSingleDocument{Value: jsonMap}, storage.IndentOption("   "))
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
}
