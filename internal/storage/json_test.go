package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var jsonBytes = []byte(`{
  "name": "Tom"
}`)
var jsonMap = map[string]interface{}{
	"name": "Tom",
}

func TestJSONParser_FromBytes(t *testing.T) {
	got, err := (&storage.JSONParser{}).FromBytes(jsonBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(jsonMap, got) {
		t.Errorf("expected %v, got %v", jsonMap, got)
	}
}

func TestJSONParser_ToBytes(t *testing.T) {
	got, err := (&storage.JSONParser{}).ToBytes(jsonMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(jsonBytes, got) {
		t.Errorf("expected %v, got %v", string(jsonBytes), string(got))
	}
}
