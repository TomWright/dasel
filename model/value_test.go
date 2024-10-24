package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestType_String(t *testing.T) {
	run := func(ty model.Type, exp string) func(*testing.T) {
		return func(t *testing.T) {
			got := ty.String()
			if got != exp {
				t.Errorf("expected %s, got %s", exp, got)
			}
		}
	}
	t.Run("string", run(model.TypeString, "string"))
	t.Run("int", run(model.TypeInt, "int"))
	t.Run("float", run(model.TypeFloat, "float"))
	t.Run("bool", run(model.TypeBool, "bool"))
	t.Run("map", run(model.TypeMap, "map"))
	t.Run("slice", run(model.TypeSlice, "array"))
	t.Run("slice", run(model.TypeUnknown, "unknown"))
	t.Run("slice", run(model.TypeNull, "null"))
}

func TestValue_Len(t *testing.T) {
	run := func(v *model.Value, exp int) func(*testing.T) {
		return func(t *testing.T) {
			got, err := v.Len()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if got != exp {
				t.Errorf("expected %d, got %d", exp, got)
			}
		}
	}
	t.Run("string", func(t *testing.T) {
		t.Run("empty", run(model.NewStringValue(""), 0))
		t.Run("non-empty", run(model.NewStringValue("hello"), 5))
	})
	t.Run("slice", func(t *testing.T) {
		t.Run("empty", run(model.NewSliceValue(), 0))
		t.Run("non-empty", run(model.NewValue([]any{1, 2, 3}), 3))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("empty", run(model.NewMapValue(), 0))
		t.Run("non-empty", run(model.NewValue(map[string]any{"one": 1, "two": 2, "three": 3}), 3))
	})
}
