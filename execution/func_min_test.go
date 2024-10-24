package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncMin(t *testing.T) {
	t.Run("int", testCase{
		s:   `min(1, 2, 3)`,
		out: model.NewIntValue(1),
	}.run)
	t.Run("float", testCase{
		s:   `min(1f, 2.5, 3.5)`,
		out: model.NewFloatValue(1),
	}.run)
	t.Run("mixed", testCase{
		s:   `min(1, 2f)`,
		out: model.NewIntValue(1),
	}.run)
}
