package model_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

type goValueTestCase struct {
	in   *model.Value
	inFn func() *model.Value
	exp  any
}

func (tc goValueTestCase) run(t *testing.T) {
	if tc.inFn != nil {
		tc.in = tc.inFn()
	}
	out, err := tc.in.GoValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cmp.Equal(tc.exp, out) {
		t.Errorf("unexpected result: %s", cmp.Diff(tc.exp, out))
	}
}

func TestValue_GoValue(t *testing.T) {
	t.Run("null", goValueTestCase{
		in:  model.NewNullValue(),
		exp: nil,
	}.run)
	t.Run("int", goValueTestCase{
		in:  model.NewIntValue(42),
		exp: int64(42),
	}.run)
	t.Run("float", goValueTestCase{
		in:  model.NewFloatValue(3.14),
		exp: 3.14,
	}.run)
	t.Run("string", goValueTestCase{
		in:  model.NewStringValue("hello"),
		exp: "hello",
	}.run)
	t.Run("bool", goValueTestCase{
		in:  model.NewBoolValue(true),
		exp: true,
	}.run)
	t.Run("slice", goValueTestCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewIntValue(1))
			_ = s.Append(model.NewIntValue(2))
			_ = s.Append(model.NewIntValue(3))
			return s
		},
		exp: []any{int64(1), int64(2), int64(3)},
	}.run)
	t.Run("map", goValueTestCase{
		inFn: func() *model.Value {
			m := model.NewMapValue()
			_ = m.SetMapKey("a", model.NewStringValue("apple"))
			return m
		},
		exp: map[string]any{
			"a": "apple",
		},
	}.run)
}
