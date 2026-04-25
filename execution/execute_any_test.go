package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestAnyExpr(t *testing.T) {
	t.Run("some match", testCase{
		s:   `[1, 2, 3].any($this > 2)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("none match", testCase{
		s:   `[1, 2, 3].any($this > 5)`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("all match", testCase{
		s:   `[1, 2, 3].any($this > 0)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("with equality", testCase{
		s:   `["a", "b", "c"].any($this == "b")`,
		out: model.NewBoolValue(true),
	}.run)
}
