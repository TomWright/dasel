package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncToUpper(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `"Hello World".toUpper()`,
		out: model.NewStringValue("HELLO WORLD"),
	}.run)
	t.Run("arg input", testCase{
		s:   `toUpper("Hello World")`,
		out: model.NewStringValue("HELLO WORLD"),
	}.run)
	t.Run("already uppercase", testCase{
		s:   `"ABC".toUpper()`,
		out: model.NewStringValue("ABC"),
	}.run)
	t.Run("all lowercase", testCase{
		s:   `"abc".toUpper()`,
		out: model.NewStringValue("ABC"),
	}.run)
}
