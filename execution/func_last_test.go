package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncLast(t *testing.T) {
	t.Run("int array", testCase{
		s:   `last([1, 2, 3])`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("string array", testCase{
		s:   `last(["a", "b", "c"])`,
		out: model.NewStringValue("c"),
	}.run)
	t.Run("single element", testCase{
		s:   `last([42])`,
		out: model.NewIntValue(42),
	}.run)
	t.Run("empty array", testCase{
		s:   `last([])`,
		out: model.NewNullValue(),
	}.run)
	t.Run("chained", testCase{
		s:   `[1, 2, 3].last()`,
		out: model.NewIntValue(3),
	}.run)
}
