package command

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPut_Object(t *testing.T) {
	t.Run("SingleFailingWriter", func(t *testing.T) {
		err := runPutObjectCommand(putObjectOpts{
			Parser:      "json",
			Selector:    ".[0]",
			Reader:      strings.NewReader(`[{"name": "Tom"}]`),
			InputValues: []string{"name=Frank"},
			InputTypes:  []string{"string"},
			Writer:      &failingWriter{},
		}, nil)

		if err == nil || !errors.Is(err, errFailingWriterErr) {
			t.Errorf("expected error %v, got %v", errFailingWriterErr, err)
			return
		}
	})
	t.Run("MultiFailingWriter", func(t *testing.T) {
		err := runPutObjectCommand(putObjectOpts{
			Parser:      "json",
			Selector:    ".[*]",
			Reader:      strings.NewReader(`[{"name": "Tom"}]`),
			InputValues: []string{"name=Frank"},
			InputTypes:  []string{"string"},
			Writer:      &failingWriter{},
			Multi:       true,
		}, nil)

		if err == nil || !errors.Is(err, errFailingWriterErr) {
			t.Errorf("expected error %v, got %v", errFailingWriterErr, err)
			return
		}
	})
}

func TestGetMapFromTypes(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		exp := map[string]interface{}{
			"name":   "Tom",
			"age":    int64(27),
			"active": true,
		}
		got, err := getMapFromTypesValues([]string{"string", "int", "bool"}, []string{"name=Tom", "age=27", "active=true"})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
	t.Run("InvalidTypes", func(t *testing.T) {
		_, err := getMapFromTypesValues([]string{}, []string{"name=Tom"})
		if exp := "exactly 1 types are required, got 0"; err == nil || exp != err.Error() {
			t.Errorf("expected %v, got %v", exp, err)
			return
		}
	})
	t.Run("InvalidValue", func(t *testing.T) {
		_, err := getMapFromTypesValues([]string{"int"}, []string{"x=asd"})
		if exp := "could not parse value [x]: could not parse int [asd]: strconv.ParseInt: parsing \"asd\": invalid syntax"; err == nil || exp != err.Error() {
			t.Errorf("expected %v, got %v", exp, err)
			return
		}
	})
}
