package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
)

func TestFuncKeys(t *testing.T) {
	t.Run("returns slice indices", testCase{
		s: `[1,2,3,4,5].keys()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for i := 0; i < 5; i++ {
				if err := r.Append(model.NewIntValue(int64(i))); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
	t.Run("returns map keys", testCase{
		s: `{"a": 3, "b": 4, "c": 5}.keys()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()
			for _, key := range []string{"a", "b", "c"} {
				if err := r.Append(model.NewStringValue(key)); err != nil {
					panic(err)
				}
			}
			return r
		},
	}.run)
}
