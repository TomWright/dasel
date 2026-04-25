package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncFlatten(t *testing.T) {
	t.Run("nested arrays", testCase{
		s: `flatten([[1, 2], [3, 4]])`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []int64{1, 2, 3, 4} {
				if err := r.Append(model.NewIntValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("mixed nested and flat", testCase{
		s: `flatten([1, [2, 3], 4])`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []int64{1, 2, 3, 4} {
				if err := r.Append(model.NewIntValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("already flat", testCase{
		s: `flatten([1, 2, 3])`,
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
	t.Run("chained", testCase{
		s: `[[1, 2], [3, 4]].flatten()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []int64{1, 2, 3, 4} {
				if err := r.Append(model.NewIntValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
}
