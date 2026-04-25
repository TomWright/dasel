package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncAvg(t *testing.T) {
	t.Run("ints", testCase{
		s:   `avg(2, 4, 6)`,
		out: model.NewFloatValue(4.0),
	}.run)
	t.Run("floats", testCase{
		s:   `avg(1.5, 2.5, 3.0)`,
		out: model.NewFloatValue(7.0 / 3.0),
	}.run)
	t.Run("mixed", testCase{
		s:   `avg(1, 2.0, 3)`,
		out: model.NewFloatValue(2.0),
	}.run)
	t.Run("single value", testCase{
		s:   `avg(5)`,
		out: model.NewFloatValue(5.0),
	}.run)
}
