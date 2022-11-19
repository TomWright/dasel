package dasel

import (
	"testing"
)

func TestNotFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"not()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "not",
			Args:     []string{},
		}),
	)

	t.Run(
		"Single Equal",
		selectTest(
			"name.all().not(equal(key(),first))",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				false,
				true,
			},
		),
	)

	t.Run(
		"Not Banned",
		selectTest(
			"all().filter(not(equal(banned,true))).name",
			[]map[string]interface{}{
				{
					"name":   "Tom",
					"banned": true,
				},
				{
					"name":   "Jess",
					"banned": false,
				},
			},
			[]interface{}{
				"Jess",
			},
		),
	)
}
