package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestArray(t *testing.T) {
	inSlice := func() *model.Value {
		s := model.NewSliceValue()
		if err := s.Append(model.NewIntValue(1)); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := s.Append(model.NewIntValue(2)); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := s.Append(model.NewIntValue(3)); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		return s
	}
	inMap := func() *model.Value {
		m := model.NewMapValue()
		if err := m.SetMapKey("numbers", inSlice()); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		return m
	}

	runArrayTests := func(in func() *model.Value, prefix string) func(t *testing.T) {
		return func(t *testing.T) {
			t.Run("1:2", testCase{
				s:    prefix + `[1:2]`,
				inFn: in,
				outFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewIntValue(2)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(3)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					return res
				},
			}.run)
			t.Run("1:0", testCase{
				s:    prefix + `[1:0]`,
				inFn: in,
				outFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewIntValue(2)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(1)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					return res
				},
			}.run)
			t.Run("1:", testCase{
				s:    prefix + `[1:]`,
				inFn: in,
				outFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewIntValue(2)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(3)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					return res
				},
			}.run)
			t.Run(":1", testCase{
				s:    prefix + `[:1]`,
				inFn: in,
				outFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewIntValue(1)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(2)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					return res
				},
			}.run)
			t.Run("reverse", testCase{
				s:    prefix + `[len($this)-1:0]`,
				inFn: in,
				outFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewIntValue(3)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(2)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					if err := res.Append(model.NewIntValue(1)); err != nil {
						t.Fatalf("unexpected error: %s", err)
					}
					return res
				},
			}.run)
		}
	}

	t.Run("direct to slice", runArrayTests(inSlice, "$this"))
	t.Run("property to slice", runArrayTests(inMap, "numbers"))
}
