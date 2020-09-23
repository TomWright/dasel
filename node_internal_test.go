package dasel

import (
	"reflect"
	"testing"
)

func TestSafeIsNil(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		if exp, got := false, safeIsNil(reflect.ValueOf("")); exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	})
}
