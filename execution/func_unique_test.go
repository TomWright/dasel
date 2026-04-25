package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncUnique(t *testing.T) {
	t.Run("ints with duplicates", testCase{
		s: `unique([1, 2, 2, 3, 1])`,
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
	t.Run("strings with duplicates", testCase{
		s: `unique(["a", "b", "a", "c"])`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []string{"a", "b", "c"} {
				if err := r.Append(model.NewStringValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("already unique", testCase{
		s: `unique([1, 2, 3])`,
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
		s: `[1, 1, 2].unique()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, v := range []int64{1, 2} {
				if err := r.Append(model.NewIntValue(v)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
}
