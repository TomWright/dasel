package dasel

import (
	"testing"
)

func TestLastFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"last(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "last",
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
		"Last",
		selectTest(
			"colours.last()",
			original,
			[]interface{}{
				"blue",
			},
		),
	)
}
