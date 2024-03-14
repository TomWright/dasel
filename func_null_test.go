package dasel

import (
	"testing"
)

func TestNullFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"null(1)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "null",
			Args:     []string{"1"},
		}),
	)

	original := map[string]interface{}{}

	t.Run(
		"Null",
		selectTest(
			"null()",
			original,
			[]interface{}{
				nil,
			},
		),
	)
}
