package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

func TestFuncSum(t *testing.T) {
	t.Run("int", testCase{
		s:   `sum(1, 2, 3)`,
		out: model.NewIntValue(6),
	}.run)
	t.Run("float", testCase{
		s:   `sum(1.1, 2.2, 3.3)`,
		out: model.NewFloatValue(6.6),
	}.run)
	t.Run("negative int", testCase{
		s:   `sum(-1, -2, -3)`,
		out: model.NewIntValue(-6),
	}.run)
	t.Run("negative float", testCase{
		s:   `sum(-1.1, -2.2, -3.3)`,
		out: model.NewFloatValue(-6.6),
	}.run)
	t.Run("using int and float together returns float", testCase{
		s:   `sum(1, 1.1)`,
		out: model.NewFloatValue(2.1),
	}.run)
}
