package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncAbs(t *testing.T) {
	t.Run("positive int", testCase{
		s:   `abs(5)`,
		out: model.NewIntValue(5),
	}.run)
	t.Run("negative int", testCase{
		s:   `abs(-5)`,
		out: model.NewIntValue(5),
	}.run)
	t.Run("zero", testCase{
		s:   `abs(0)`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("positive float", testCase{
		s:   `abs(3.14)`,
		out: model.NewFloatValue(3.14),
	}.run)
	t.Run("negative float", testCase{
		s:   `abs(-3.14)`,
		out: model.NewFloatValue(3.14),
	}.run)
	t.Run("chained int", testCase{
		s:   `(-5).abs()`,
		out: model.NewIntValue(5),
	}.run)
}
