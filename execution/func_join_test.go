package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
)

func TestFuncJoin(t *testing.T) {
	t.Run("chained input", testCase{
		s:   `["a","b","c"].join(",")`,
		out: model.NewStringValue("a,b,c"),
	}.run)
	t.Run("vararg input", testCase{
		s:   `join(",", "a", "b", "c")`,
		out: model.NewStringValue("a,b,c"),
	}.run)
	t.Run("array input", testCase{
		s:   `join(",", ["a", "b", "c"])`,
		out: model.NewStringValue("a,b,c"),
	}.run)
}
