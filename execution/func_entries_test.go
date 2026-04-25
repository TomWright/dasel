package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncEntries(t *testing.T) {
	t.Run("map to entries", testCase{
		s: `{"a": 1, "b": 2}.entries()`,
		outFn: func() *model.Value {
			r := model.NewSliceValue()

			e1 := model.NewMapValue()
			_ = e1.SetMapKey("key", model.NewStringValue("a"))
			_ = e1.SetMapKey("value", model.NewIntValue(1))
			_ = r.Append(e1)

			e2 := model.NewMapValue()
			_ = e2.SetMapKey("key", model.NewStringValue("b"))
			_ = e2.SetMapKey("value", model.NewIntValue(2))
			_ = r.Append(e2)

			return r
		},
	}.run)
}

func TestFuncFromEntries(t *testing.T) {
	t.Run("entries to map", testCase{
		s: `[{"key": "a", "value": 1}, {"key": "b", "value": 2}].fromEntries()`,
		outFn: func() *model.Value {
			r := model.NewMapValue()
			_ = r.SetMapKey("a", model.NewIntValue(1))
			_ = r.SetMapKey("b", model.NewIntValue(2))
			return r
		},
	}.run)
	t.Run("roundtrip", testCase{
		s: `{"x": 10, "y": 20}.entries().fromEntries()`,
		outFn: func() *model.Value {
			r := model.NewMapValue()
			_ = r.SetMapKey("x", model.NewIntValue(10))
			_ = r.SetMapKey("y", model.NewIntValue(20))
			return r
		},
	}.run)
}
