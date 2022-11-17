package query

import (
	"testing"
)

func TestParentFunc(t *testing.T) {
	t.Run(
		"SimpleParent",
		selectTest(
			"name.first.parent()",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"MultiParent",
		selectTest(
			"name.all().parent()",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)
}
