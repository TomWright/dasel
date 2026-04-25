package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncValues(t *testing.T) {
	t.Run("map values", testCase{
		s: `{"a": 1, "b": 2, "c": 3}.values()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []int64{1, 2, 3} {
				if err := r.Append(model.NewIntValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("single entry", testCase{
		s: `{"x": 42}.values()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			if err := r.Append(model.NewIntValue(42)); err != nil {
				panic(err)
			}
			return r
		},
	}.run)
}
