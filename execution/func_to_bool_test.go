package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncToBool(t *testing.T) {
	t.Run("true bool", testCase{
		s:   `toBool(true)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("false bool", testCase{
		s:   `toBool(false)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("string true", testCase{
		s:   `toBool("true")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("string false", testCase{
		s:   `toBool("false")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("string yes", testCase{
		s:   `toBool("yes")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("string no", testCase{
		s:   `toBool("no")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("int 1", testCase{
		s:   `toBool(1)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("int 0", testCase{
		s:   `toBool(0)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("null", testCase{
		s:   `toBool(null)`,
		out: model.NewBoolValue(false),
	}.run)
}
