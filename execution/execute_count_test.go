package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestCountExpr(t *testing.T) {
	t.Run("some match", testCase{
		s:   `[1, 2, 3, 4, 5].count($this > 3)`,
		out: model.NewIntValue(2),
	}.run)
	t.Run("all match", testCase{
		s:   `[1, 2, 3].count($this > 0)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].count($this > 5)`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "b", "a", "c", "a"].count($this == "a")`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("empty array", testCase{
		s:   `[].count($this > 0)`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("boolean literal true", testCase{
		s:   `[1, 2, 3].count(true)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("boolean literal false", testCase{
		s:   `[1, 2, 3].count(false)`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("with nested property", testCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("status", "active")))
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("status", "inactive")))
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("status", "active")))
			return s
		},
		s:   `count(status == "active")`,
		out: model.NewIntValue(2),
	}.run)
	t.Run("chained after filter", testCase{
		s:   `[1, 2, 3, 4, 5].filter($this > 2).count($this > 4)`,
		out: model.NewIntValue(1),
	}.run)
}
