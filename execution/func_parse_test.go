package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
	"testing"
)

func TestFuncParse(t *testing.T) {
	t.Run("json", testCase{
		s:   `parse('json', '{"foo":"bar"}').foo`,
		out: model.NewStringValue("bar"),
	}.run)
}
