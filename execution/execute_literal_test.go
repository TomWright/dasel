package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestLiteral(t *testing.T) {
	t.Run("string", testCase{
		s:   `"hello"`,
		out: model.NewStringValue("hello"),
	}.run)
	t.Run("int", testCase{
		s:   `123`,
		out: model.NewIntValue(123),
	}.run)
	t.Run("float", testCase{
		s:   `123.4`,
		out: model.NewFloatValue(123.4),
	}.run)
	t.Run("true", testCase{
		s:   `true`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("false", testCase{
		s:   `false`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("empty array", testCase{
		s: `[]`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			return r
		},
	}.run)
	t.Run("array with one element", testCase{
		s: `[1]`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
	}.run)
	t.Run("array with many elements", testCase{
		s: `[1, 2.2, "foo", true, [1, 2, 3]]`,
		outFn: func() *model.Value {
			nested := model.NewSliceValue()
			if err := nested.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := nested.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := nested.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			r := model.NewSliceValue()
			if err := r.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewFloatValue(2.2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewStringValue("foo")); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewBoolValue(true)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(nested); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
	}.run)
	t.Run("array with expressions", testCase{
		s: `[1 + 1, 2f - 2, "foo" + "bar", true || false, [1 + 1, 2 * 2, 3 / 3]]`,
		outFn: func() *model.Value {
			nested := model.NewSliceValue()
			if err := nested.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := nested.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := nested.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			r := model.NewSliceValue()
			if err := r.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewFloatValue(0)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewStringValue("foobar")); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(model.NewBoolValue(true)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := r.Append(nested); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return r
		},
	}.run)
}
