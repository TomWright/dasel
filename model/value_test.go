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
	t.Run("unknown", run(model.TypeUnknown, "unknown"))
	t.Run("null", run(model.TypeNull, "null"))
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

func TestValue_IsScalar(t *testing.T) {
	run := func(v *model.Value, exp bool) func(*testing.T) {
		return func(t *testing.T) {
			got := v.IsScalar()
			if got != exp {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("string", run(model.NewStringValue("foo"), true))
	t.Run("bool", run(model.NewBoolValue(true), true))
	t.Run("int", run(model.NewIntValue(1), true))
	t.Run("float", run(model.NewFloatValue(1.0), true))
	t.Run("null", run(model.NewNullValue(), true))
	t.Run("map", run(model.NewMapValue(), false))
	t.Run("slice", run(model.NewSliceValue(), false))

	t.Run("nested", func(t *testing.T) {
		t.Run("nested string", run(model.NewNestedValue(model.NewStringValue("foo")), true))
		t.Run("nested bool", run(model.NewNestedValue(model.NewBoolValue(true)), true))
		t.Run("nested int", run(model.NewNestedValue(model.NewIntValue(1)), true))
		t.Run("nested float", run(model.NewNestedValue(model.NewFloatValue(1.0)), true))
		t.Run("nested null", run(model.NewNestedValue(model.NewNullValue()), true))
		t.Run("nested map", run(model.NewNestedValue(model.NewMapValue()), false))
		t.Run("nested slice", run(model.NewNestedValue(model.NewSliceValue()), false))

		t.Run("double nested string", run(model.NewNestedValue(model.NewNestedValue(model.NewStringValue("foo"))), true))
	})
}
