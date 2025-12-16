package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

func TestFuncToInt(t *testing.T) {
	t.Run("string", testCase{
		s:   `toInt("2")`,
		out: model.NewIntValue(2),
	}.run)
	t.Run("int", testCase{
		s:   `toInt(1)`,
		out: model.NewIntValue(1),
	}.run)
	t.Run("float", testCase{
		s:   `toInt(1.1)`,
		out: model.NewIntValue(1),
	}.run)
	t.Run("bool", testCase{
		s:   `toInt(true)`,
		out: model.NewIntValue(1),
	}.run)
}
