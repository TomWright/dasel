package dasel

import (
	"testing"
)

func TestMapOfFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"mapOf()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "mapOf",
			Args:     []string{},
		}),
	)

	t.Run(
		"Single Equal",
		selectTest(
			"mapOf(firstName,name.first,lastName,name.last)",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				map[string]interface{}{
					"firstName": "Tom",
					"lastName":  "Wright",
				},
			},
		),
	)
}
