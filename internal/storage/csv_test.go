package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"testing"
)

var csvBytes = []byte(`id,name
1,Tom
2,Jim

`)
var csvMap = []map[string]string{
	{
		"id":   "1",
		"name": "Tom",
	},
	{
		"id":   "2",
		"name": "Jim",
	},
}

func TestCSVParser_FromBytes(t *testing.T) {
	got, err := (&storage.CSVParser{}).FromBytes(csvBytes)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(csvMap, got) {
		t.Errorf("expected %v, got %v", csvMap, got)
	}
}

func TestCSVParser_FromBytes_Error(t *testing.T) {
	_, err := (&storage.CSVParser{}).FromBytes(nil)
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.CSVParser{}).FromBytes([]byte(`a,b
a,b,c`))
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
	_, err = (&storage.CSVParser{}).FromBytes([]byte(`a,b,c
a,b`))
	if err == nil {
		t.Errorf("expected error but got none")
		return
	}
}

func TestCSVParser_ToBytes(t *testing.T) {
	got, err := (&storage.CSVParser{}).ToBytes(csvMap)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(csvBytes, got) {
		t.Errorf("expected %v, got %v", string(csvBytes), string(got))
	}
}
