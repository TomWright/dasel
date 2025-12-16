package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

func TestFuncToFloat(t *testing.T) {
	t.Run("string", testCase{
		s:   `toFloat("1.1")`,
		out: model.NewFloatValue(1.1),
	}.run)
	t.Run("int", testCase{
		s:   `toFloat(1)`,
		out: model.NewFloatValue(1),
	}.run)
	t.Run("float", testCase{
		s:   `toFloat(1.1)`,
		out: model.NewFloatValue(1.1),
	}.run)
	t.Run("bool", testCase{
		s:   `toFloat(true)`,
		out: model.NewFloatValue(1),
	}.run)
}
