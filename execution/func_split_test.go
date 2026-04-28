package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncSplit(t *testing.T) {
	t.Run("chained input", testCase{
		s: `"a,b,c".split(",")`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, s := range []string{"a", "b", "c"} {
				if err := r.Append(model.NewStringValue(s)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("arg input", testCase{
		s: `split(",", "a,b,c")`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, s := range []string{"a", "b", "c"} {
				if err := r.Append(model.NewStringValue(s)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("no match", testCase{
		s: `"abc".split(",")`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			if err := r.Append(model.NewStringValue("abc")); err != nil {
				panic(err)
			}
			return r
		},
	}.run)
	t.Run("empty separator", testCase{
		s: `"abc".split("")`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, s := range []string{"a", "b", "c"} {
				if err := r.Append(model.NewStringValue(s)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
}
