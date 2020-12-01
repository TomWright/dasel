package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var xmlBytes = []byte(`<user>
  <name>Tom</name>
</user>
`)
var xmlMap = map[string]interface{}{
	"user": map[string]interface{}{
		"name": "Tom",
	},
}

func TestXMLParser_FromBytes(t *testing.T) {
	got, err := (&storage.XMLParser{}).FromBytes(xmlBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(&storage.BasicSingleDocument{Value: xmlMap}, got) {
		t.Errorf("expected %v, got %v", xmlMap, got)
	}
}

func TestXMLParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.XMLParser{}).FromBytes(nil)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.XMLParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestXMLParser_ToBytes(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes(xmlMap)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(xmlBytes, got) {
			t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
		}
	})
	t.Run("SingleDocument", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: xmlMap})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(xmlBytes, got) {
			t.Errorf("expected %v, got %v", string(xmlBytes), string(got))
		}
	})
	t.Run("MultiDocument", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{xmlMap, xmlMap}})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := append([]byte{}, xmlBytes...)
		exp = append(exp, xmlBytes...)
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", string(exp), string(got))
		}
	})
	t.Run("DefaultValue", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes("asd")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := []byte(`asd
`)
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", string(exp), string(got))
		}
	})
	t.Run("SingleDocumentValue", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicSingleDocument{Value: "asd"})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := []byte(`asd
`)
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", string(exp), string(got))
		}
	})
	t.Run("MultiDocumentValue", func(t *testing.T) {
		got, err := (&storage.XMLParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{"asd", "123"}})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := []byte(`asd
123
`)
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", string(exp), string(got))
		}
	})
}
