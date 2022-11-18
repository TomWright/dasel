package dasel

import (
	"testing"
)

func TestIndexFunc(t *testing.T) {
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
		"Index",
		selectTest(
			"colours.index(1)",
			original,
			[]interface{}{
				"green",
			},
		),
	)

	t.Run(
		"IndexShorthand",
		selectTest(
			"colours.[1]",
			original,
			[]interface{}{
				"green",
			},
		),
	)
}
