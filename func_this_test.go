package dasel

import (
	"testing"
)

func TestThisFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"this(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "this",
			Args:     []string{"x"},
		}),
	)
	t.Run(
		"SimpleThis",
		selectTest(
			"name.this().first",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				"Tom",
			},
		),
	)
	t.Run(
		"BlankSelectorThis",
		selectTest(
			".name.first",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				"Tom",
			},
		),
	)
}
