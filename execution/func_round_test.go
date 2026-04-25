package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncRound(t *testing.T) {
	t.Run("rounds up", testCase{
		s:   `round(3.7)`,
		out: model.NewIntValue(4),
	}.run)
	t.Run("rounds down", testCase{
		s:   `round(3.2)`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("half rounds up", testCase{
		s:   `round(3.5)`,
		out: model.NewIntValue(4),
	}.run)
	t.Run("negative", testCase{
		s:   `round(-3.7)`,
		out: model.NewIntValue(-4),
	}.run)
	t.Run("int passthrough", testCase{
		s:   `round(5)`,
		out: model.NewIntValue(5),
	}.run)
	t.Run("chained", testCase{
		s:   `(3.7).round()`,
		out: model.NewIntValue(4),
	}.run)
}
