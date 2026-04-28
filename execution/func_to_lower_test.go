package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncToLower(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `"Hello World".toLower()`,
		out: model.NewStringValue("hello world"),
	}.run)
	t.Run("arg input", testCase{
		s:   `toLower("Hello World")`,
		out: model.NewStringValue("hello world"),
	}.run)
	t.Run("already lowercase", testCase{
		s:   `"abc".toLower()`,
		out: model.NewStringValue("abc"),
	}.run)
	t.Run("all uppercase", testCase{
		s:   `"ABC".toLower()`,
		out: model.NewStringValue("abc"),
	}.run)
}
