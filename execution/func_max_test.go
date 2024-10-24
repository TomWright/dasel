package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncMax(t *testing.T) {
	t.Run("int", testCase{
		s:   `max(1, 2, 3)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("float", testCase{
		s:   `max(1f, 2.5, 3.5)`,
		out: model.NewFloatValue(3.5),
	}.run)
	t.Run("mixed", testCase{
		s:   `max(1, 2f)`,
		out: model.NewFloatValue(2),
	}.run)
}
