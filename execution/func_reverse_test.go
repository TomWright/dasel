package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncReverse(t *testing.T) {
	t.Run("array", testCase{
		s: `reverse([1, 2, 3])`,
		outFn: func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewIntValue(3)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := res.Append(model.NewIntValue(2)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := res.Append(model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return res
		},
	}.run)

	t.Run("string", testCase{
		s:   `reverse("hello")`,
		out: model.NewStringValue("olleh"),
	}.run)
}
