package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncTypeOf(t *testing.T) {
	t.Run("string", testCase{
		s:   `typeOf("hello")`,
		out: model.NewStringValue("string"),
	}.run)
	t.Run("int", testCase{
		s:   `typeOf(123)`,
		out: model.NewStringValue("int"),
	}.run)
	t.Run("float", testCase{
		s:   `typeOf(12.3)`,
		out: model.NewStringValue("float"),
	}.run)
	t.Run("bool", testCase{
		s:   `typeOf(true)`,
		out: model.NewStringValue("bool"),
	}.run)
	t.Run("array", testCase{
		s:   `typeOf([])`,
		out: model.NewStringValue("array"),
	}.run)
	t.Run("map", testCase{
		s:   `typeOf({})`,
		out: model.NewStringValue("map"),
	}.run)
	t.Run("null", testCase{
		s:   `typeOf(null)`,
		out: model.NewStringValue("null"),
	}.run)
}
