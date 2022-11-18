package dasel

import (
	"testing"
)

func TestLastFunc(t *testing.T) {
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
