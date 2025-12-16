package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
	"testing"
)

func TestFuncLen(t *testing.T) {
	t.Run("array", testCase{
		s:   `len([1,2,3])`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("object", testCase{
		s:   `len({"foo":1,"bar":2,"baz":3})`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("string", testCase{
		s:   `len("hello")`,
		out: model.NewIntValue(5),
	}.run)
}
