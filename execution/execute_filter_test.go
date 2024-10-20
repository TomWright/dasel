package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFilter(t *testing.T) {
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
	t.Run("all true", testCase{
		inFn: inSlice,
		s:    "filter(true)",
		outFn: func() *model.Value {
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
		},
	}.run)
	t.Run("all false", testCase{
		inFn: inSlice,
		s:    "filter(false)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			return s
		},
	}.run)
	t.Run("equal 2", testCase{
		inFn: inSlice,
		s:    "filter($this == 2)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			if err := s.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			return s
		},
	}.run)
	t.Run("not equal 2", testCase{
		inFn: inSlice,
		s:    "filter($this != 2)",
		outFn: func() *model.Value {
			s := model.NewSliceValue()
			if err := s.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if err := s.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			return s
		},
	}.run)
}
