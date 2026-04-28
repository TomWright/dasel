package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestAllExpr(t *testing.T) {
	t.Run("all match", testCase{
		s:   `[1, 2, 3].all($this > 0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("some match", testCase{
		s:   `[1, 2, 3].all($this > 1)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].all($this > 5)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "a", "a"].all($this == "a")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("empty array", testCase{
		s:   `[].all($this > 0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("boolean literal true", testCase{
		s:   `[1, 2, 3].all(true)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("boolean literal false", testCase{
		s:   `[1, 2, 3].all(false)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("with negation", testCase{
		s:   `[1, 2, 3].all(!($this > 5))`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("with nested property", testCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("active", true)))
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("active", true)))
			return s
		},
		s:   `all(active == true)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("chained after filter", testCase{
		s:   `[1, 2, 3, 4, 5].filter($this > 2).all($this > 2)`,
		out: model.NewBoolValue(true),
	}.run)
}
