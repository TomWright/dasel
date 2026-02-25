package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestSpread(t *testing.T) {
	t.Run("build new array", testCase{
		s: "[[1,2,3]..., 4]",
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
			if err := s.Append(model.NewIntValue(4)); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			return s
		},
	}.run)
}
