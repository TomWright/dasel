package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	_ "github.com/tomwright/dasel/v3/parsing/json"
)

func TestFuncGet(t *testing.T) {
	t.Run("returns array element", testCase{
		s:   `[1,2,3,4,5].get(3)`,
		out: model.NewIntValue(4),
	}.run)
	t.Run("returns map key", testCase{
		s:   `{"a": 3, "b": 4, "c": 5}.get("b")`,
		out: model.NewIntValue(4),
	}.run)
	t.Run("coalesce with invalid map accessor", testCase{
		s:   `{}.get("a") ?? "missing"`,
		out: model.NewStringValue("missing"),
	}.run)
	t.Run("returns null when string accessor used on slice", testCase{
		s:   `[].get(0) ?? "missing"`,
		out: model.NewStringValue("missing"),
	}.run)
}
