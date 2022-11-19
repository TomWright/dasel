package dasel

import (
	"testing"
)

func TestFirstFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"first(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "first",
			Args:     []string{"x"},
		}),
	)

	original := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tom",
			"last":  "Wright",
		},
		"colours": []interface{}{
			"red", "green", "blue",
		},
	}

	t.Run(
		"First",
		selectTest(
			"colours.first()",
			original,
			[]interface{}{
				"red",
			},
		),
	)
}
