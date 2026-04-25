package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncIndexOf(t *testing.T) {
	t.Run("found", testCase{
		s:   `"hello world".indexOf("world")`,
		out: model.NewIntValue(6),
	}.run)
	t.Run("not found", testCase{
		s:   `"hello world".indexOf("xyz")`,
		out: model.NewIntValue(-1),
	}.run)
	t.Run("at start", testCase{
		s:   `"hello".indexOf("hello")`,
		out: model.NewIntValue(0),
	}.run)
	t.Run("arg input", testCase{
		s:   `indexOf("hello world", "world")`,
		out: model.NewIntValue(6),
	}.run)
	t.Run("empty substring", testCase{
		s:   `"hello".indexOf("")`,
		out: model.NewIntValue(0),
	}.run)
}
