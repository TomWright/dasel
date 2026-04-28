package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestAnyExpr(t *testing.T) {
	t.Run("some match", testCase{
		s:   `[1, 2, 3].any($this > 2)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].any($this > 5)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("all match", testCase{
		s:   `[1, 2, 3].any($this > 0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "b", "c"].any($this == "b")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("empty array", testCase{
		s:   `[].any($this > 0)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("boolean literal true", testCase{
		s:   `[1, 2, 3].any(true)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("boolean literal false", testCase{
		s:   `[1, 2, 3].any(false)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("with negation", testCase{
		s:   `[1, 2, 3].any(!($this > 2))`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("with nested property", testCase{
		inFn: func() *model.Value {
			s := model.NewSliceValue()
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("age", int64(20))))
			_ = s.Append(model.NewValue(orderedmap.NewMap().Set("age", int64(30))))
			return s
		},
		s:   `any(age > 25)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("chained after filter", testCase{
		s:   `[1, 2, 3, 4, 5].filter($this > 2).any($this > 4)`,
		out: model.NewBoolValue(true),
	}.run)
}
