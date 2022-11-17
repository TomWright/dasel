package query

import (
	"testing"
)

func TestEqualFunc(t *testing.T) {

	t.Run(
		"Single Equal",
		selectTest(
			"name.all().equal(key(),first)",
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
		"Multi Equal",
		selectTest(
			"name.all().equal(key(),first,key(),last)",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				"Tom",
				"Wright",
			},
		),
	)
	t.Run(
		"Multi Equal Into Single Equal",
		selectTest(
			"name.all().equal(key(),first,key(),last).equal(key(),first)",
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
