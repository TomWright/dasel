package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncFirst(t *testing.T) {
	t.Run("int array", testCase{
		s:   `first([1, 2, 3])`,
		out: model.NewIntValue(1),
	}.run)
	t.Run("string array", testCase{
		s:   `first(["a", "b", "c"])`,
		out: model.NewStringValue("a"),
	}.run)
	t.Run("single element", testCase{
		s:   `first([42])`,
		out: model.NewIntValue(42),
	}.run)
	t.Run("empty array", testCase{
		s:   `first([])`,
		out: model.NewNullValue(),
	}.run)
	t.Run("chained", testCase{
		s:   `[1, 2, 3].first()`,
		out: model.NewIntValue(1),
	}.run)
}
