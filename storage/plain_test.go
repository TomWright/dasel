package storage_test

import (
	"errors"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
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
		gotVal, err := (&storage.PlainParser{}).ToBytes(dasel.ValueOf("asd"))
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
		gotVal, err := (&storage.PlainParser{}).ToBytes(dasel.ValueOf("asd").WithMetadata("isSingleDocument", true))
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
		val := dasel.ValueOf([]interface{}{"asd", "123"})
		daselVal := dasel.ValueOf(val).WithMetadata("isMultiDocument", true)

		gotVal, err := (&storage.PlainParser{}).ToBytes(daselVal)
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
