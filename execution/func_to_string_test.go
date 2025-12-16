package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

func TestFuncToString(t *testing.T) {
	t.Run("string", testCase{
		s:   `toString("Hello")`,
		out: model.NewStringValue("Hello"),
	}.run)
	t.Run("int", testCase{
		s:   `toString(1)`,
		out: model.NewStringValue("1"),
	}.run)
	t.Run("float", testCase{
		s:   `toString(1.1)`,
		out: model.NewStringValue("1.1"),
	}.run)
	t.Run("bool", testCase{
		s:   `toString(true)`,
		out: model.NewStringValue("true"),
	}.run)
}
