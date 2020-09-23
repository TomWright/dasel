package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var yamlBytes = []byte(`name: Tom
`)
var yamlMap = map[interface{}]interface{}{
	"name": "Tom",
}

func TestYAMLParser_FromBytes(t *testing.T) {
	got, err := (&storage.YAMLParser{}).FromBytes(yamlBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(yamlMap, got) {
		t.Errorf("expected %v, got %v", yamlMap, got)
	}
}

func TestYAMLParser_ToBytes(t *testing.T) {
	got, err := (&storage.YAMLParser{}).ToBytes(yamlMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if string(yamlBytes) != string(got) {
		t.Errorf("expected %s, got %s", yamlBytes, got)
	}
}
