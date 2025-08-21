package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncHas(t *testing.T) {
	t.Run("index in range", testCase{
		s:   `[1,2,3].has(0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("negative index", testCase{
		s:   `[1,2,3].has(-1)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("index overflow", testCase{
		s:   `[1,2,3].has(3)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("index string", testCase{
		s:   `[1,2,3].has("foo")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("has map key", testCase{
		s:   `{"x":1}.has("x")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("does not have map key", testCase{
		s:   `{"x":1}.has("y")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("does not have map index", testCase{
		s:   `{"x":1}.has(1)`,
		out: model.NewBoolValue(false),
	}.run)
}
