package storage_test

import (
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/dencoding"
	"github.com/tomwright/dasel/v2/storage"
	"reflect"
	"strings"
	"testing"
)

var yamlBytes = []byte(`name: Tom
numbers:
    - 1
    - 2
`)
var yamlMap = dencoding.NewMap().
	Set("name", "Tom").
	Set("numbers", []interface{}{
		int64(1),
		int64(2),
	})

var yamlBytesMulti = []byte(`name: Tom
---
name: Jim
`)
var yamlMapMulti = []interface{}{
	dencoding.NewMap().Set("name", "Tom"),
	dencoding.NewMap().Set("name", "Jim"),
}

func TestYAMLParser_FromBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		gotFromBytes, err := (&storage.YAMLParser{}).FromBytes(yamlBytes)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := yamlMap.KeyValues()
		got := gotFromBytes.Interface().(*dencoding.Map).KeyValues()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMultiDocument", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).FromBytes(yamlBytesMulti)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := yamlMapMulti

		if !reflect.DeepEqual(exp, got.Interface()) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		_, err := (&storage.YAMLParser{}).FromBytes([]byte(`{1:asd`))
		if err == nil || !strings.Contains(err.Error(), "could not unmarshal data") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
	t.Run("Empty", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).FromBytes([]byte(``))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(dasel.Value{}, got) {
			t.Errorf("expected %v, got %v", nil, got)
		}
	})
}

func TestYAMLParser_ToBytes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(dasel.ValueOf(yamlMap))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidSingle", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(dasel.ValueOf(yamlMap).WithMetadata("isSingleDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytes) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytes, got)
		}
	})
	t.Run("ValidSingleColourise", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(dasel.ValueOf(yamlMap).WithMetadata("isSingleDocument", true), storage.ColouriseOption(true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		expBuf, _ := storage.Colourise(string(yamlBytes), "yaml")
		exp := expBuf.Bytes()
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("ValidMulti", func(t *testing.T) {
		got, err := (&storage.YAMLParser{}).ToBytes(dasel.ValueOf(yamlMapMulti).WithMetadata("isMultiDocument", true))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if string(yamlBytesMulti) != string(got) {
			t.Errorf("expected %s, got %s", yamlBytesMulti, got)
		}
	})
}
