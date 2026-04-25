package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncFloor(t *testing.T) {
	t.Run("positive float", testCase{
		s:   `floor(3.7)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("negative float", testCase{
		s:   `floor(-3.2)`,
		out: model.NewIntValue(-4),
	}.run)
	t.Run("whole float", testCase{
		s:   `floor(5.0)`,
		out: model.NewIntValue(5),
	}.run)
	t.Run("int passthrough", testCase{
		s:   `floor(5)`,
		out: model.NewIntValue(5),
	}.run)
	t.Run("chained", testCase{
		s:   `(3.7).floor()`,
		out: model.NewIntValue(3),
	}.run)
}
