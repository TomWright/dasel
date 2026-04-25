package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestCountExpr(t *testing.T) {
	t.Run("some match", testCase{
		s:   `[1, 2, 3, 4, 5].count($this > 3)`,
		out: model.NewIntValue(2),
	}.run)
	t.Run("all match", testCase{
		s:   `[1, 2, 3].count($this > 0)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].count($this > 5)`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "b", "a", "c", "a"].count($this == "a")`,
		out: model.NewIntValue(3),
	}.run)
}
