package execution_test

import (
	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
	"testing"
)

func TestFuncStringify(t *testing.T) {
	t.Run("json with explicit value", testCase{
		s:   `stringify("json", {"a": 1})`,
		out: model.NewStringValue("{\n    \"a\": 1\n}"),
	}.run)
	t.Run("json chained", testCase{
		s:   `{"a": 1}.stringify("json")`,
		out: model.NewStringValue("{\n    \"a\": 1\n}"),
	}.run)
	t.Run("roundtrip", testCase{
		s:   `parse("json", stringify("json", {"a": 1})).a`,
		out: model.NewIntValue(1),
	}.run)
}
