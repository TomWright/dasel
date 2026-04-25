package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncEndsWith(t *testing.T) {
	t.Run("chained true", testCase{
		s:   `"hello world".endsWith("world")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("chained false", testCase{
		s:   `"hello world".endsWith("hello")`,
		out: model.NewBoolValue(false),
	}.run)
	t.Run("arg input", testCase{
		s:   `endsWith("hello world", "world")`,
		out: model.NewBoolValue(true),
	}.run)
	t.Run("empty suffix", testCase{
		s:   `"hello".endsWith("")`,
		out: model.NewBoolValue(true),
	}.run)
}
