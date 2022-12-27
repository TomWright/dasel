package dasel

import (
	"testing"
)

func TestStringFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"string()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "string",
			Args:     []string{},
		}),
	)

	original := map[string]interface{}{}

	t.Run(
		"String",
		selectTest(
			"string(x)",
			original,
			[]interface{}{
				"x",
			},
		),
	)

	t.Run(
		"Comma",
		selectTest(
			"string(\\,)",
			original,
			[]interface{}{
				",",
			},
		),
	)
}
