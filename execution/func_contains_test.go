package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
	"testing"
)

func TestFuncContains(t *testing.T) {
	t.Run("array true", testCase{
		s:   `[1,2,3,4,5].contains(3)`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("array false", testCase{
		s:   `[1,2,3,4,5].contains(6)`,
		out: model.NewBoolValue(false),
	}.run)
}
