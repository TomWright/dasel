package dasel

import (
	"testing"
)

func TestIndexFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"index()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "index",
			Args:     []string{},
		}),
	)

	t.Run("NotFound", selectTestErr(
		"[0]",
		[]interface{}{},
		&ErrIndexNotFound{
			Index: 0,
		}),
	)

	t.Run("NotFoundOnInvalidType", selectTestErr(
		"[0]",
		map[string]interface{}{},
		&ErrIndexNotFound{
			Index: 0,
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
