package storage_test

import (
	"errors"
	"github.com/tomwright/dasel/internal/storage"
	"testing"
)

func TestPlainParser_FromBytes(t *testing.T) {
	_, err := (&storage.PlainParser{}).FromBytes(nil)
	if !errors.Is(err, storage.ErrPlainParserNotImplemented) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPlainParser_ToBytes(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		gotVal, err := (&storage.PlainParser{}).ToBytes("asd")
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
`
		got := string(gotVal)
		if exp != got {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("SingleDocument", func(t *testing.T) {
		gotVal, err := (&storage.PlainParser{}).ToBytes(&storage.BasicSingleDocument{Value: "asd"})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
`
		got := string(gotVal)
		if exp != got {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
	t.Run("MultiDocument", func(t *testing.T) {
		gotVal, err := (&storage.PlainParser{}).ToBytes(&storage.BasicMultiDocument{Values: []interface{}{"asd", "123"}})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := `asd
123
`
		got := string(gotVal)
		if exp != got {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	})
}
