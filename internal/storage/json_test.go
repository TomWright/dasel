package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var jsonBytes = []byte(`{"name":"Tom"}
`)
var jsonMap = map[string]interface{}{
	"name": "Tom",
}

func TestJSONParser_FromBytes(t *testing.T) {
	got, err := (&storage.JSONParser{}).FromBytes(jsonBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(jsonMap, got.RealValue()) {
		t.Errorf("expected %v, got %v", jsonMap, got)
	}
}

func TestJSONParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.JSONParser{}).FromBytes(yamlBytes)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestJSONParser_ToBytes(t *testing.T) {
	got, err := (&storage.JSONParser{}).ToBytes(jsonMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if string(jsonBytes) != string(got) {
		t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
	}
}

var jsonBytesMulti = []byte(`
{"name": "Tom"}
{"name": "Ellis"}
`)

func TestJSONParser_FromBytes_Multi(t *testing.T) {
	got, err := (&storage.JSONParser{}).FromBytes(jsonBytesMulti)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := &storage.JSONMultiDocument{
		Values: []interface{}{
			map[string]interface{}{"name": "Tom"},
			map[string]interface{}{"name": "Ellis"},
		},
	}
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", jsonMap, got)
	}
}

var mixedJsonBytesMulti = []byte(`
{
  "name": "Tom",
  "other": true
}
{"name": "Ellis"}
`)

func TestJSONParser_FromBytes_Multi_Mixed(t *testing.T) {
	got, err := (&storage.JSONParser{}).FromBytes(mixedJsonBytesMulti)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	exp := &storage.JSONMultiDocument{
		Values: []interface{}{
			map[string]interface{}{"name": "Tom", "other": true},
			map[string]interface{}{"name": "Ellis"},
		},
	}
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", jsonMap, got)
	}
}
