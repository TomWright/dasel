package dasel

import (
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
)

func TestKeysFunc(t *testing.T) {
	testdata := map[string]any{
		"object": map[string]any{
			"c": 3, "a": 1, "b": 2,
		},
		"list":   []any{111, 222, 333},
		"string": "something",
		"dencodingMap": dencoding.NewMap().
			Set("a", 1).
			Set("b", 2).
			Set("c", 3),
	}

	t.Run(
		"root",
		selectTest(
			"keys()",
			testdata,
			[]any{[]any{"dencodingMap", "list", "object", "string"}},
		),
	)

	t.Run(
		"List",
		selectTest(
			"list.keys()",
			testdata,
			[]any{[]any{0, 1, 2}},
		),
	)

	t.Run(
		"Object",
		selectTest(
			"object.keys()",
			testdata,
			[]any{[]any{"a", "b", "c"}}, // sorted
		),
	)

	t.Run(
		"Dencoding Map",
		selectTest(
			"dencodingMap.keys()",
			testdata,
			[]any{[]any{"a", "b", "c"}}, // sorted
		),
	)

	t.Run("InvalidType",
		selectTestErr(
			"string.keys()",
			testdata,
			&ErrInvalidType{
				ExpectedTypes: []string{"slice", "array", "map"},
				CurrentType:   "string",
			},
		),
	)
}
