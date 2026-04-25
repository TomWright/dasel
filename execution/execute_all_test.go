package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestAllExpr(t *testing.T) {
	t.Run("all match", testCase{
		s:   `[1, 2, 3].all($this > 0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("some match", testCase{
		s:   `[1, 2, 3].all($this > 1)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].all($this > 5)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "a", "a"].all($this == "a")`,
		out: model.NewBoolValue(true),
	}.run)
}
