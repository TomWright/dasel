package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncStartsWith(t *testing.T) {
	t.Run("chained true", testCase{
		s:   `"hello world".startsWith("hello")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("chained false", testCase{
		s:   `"hello world".startsWith("world")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("arg input", testCase{
		s:   `startsWith("hello world", "hello")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("empty prefix", testCase{
		s:   `"hello".startsWith("")`,
		out: model.NewBoolValue(true),
	}.run)
}
