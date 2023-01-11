package dasel

import (
	"testing"
)

func TestParentFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"parent(x,x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "parent",
			Args:     []string{"x", "x"},
		}),
	)

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
		"SimpleParent2Levels",
		selectTest(
			"user.name.first.parent(2).deleted",
			map[string]interface{}{
				"user": map[string]interface{}{
					"name": map[string]interface{}{
						"first": "Tom",
						"last":  "Wright",
					},
					"deleted": false,
				},
			},
			[]interface{}{
				false,
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

	t.Run(
		"FilteredParent",
		selectTest(
			"all().flags.filter(equal(banned,false)).parent().name",
			[]map[string]interface{}{
				{
					"flags": map[string]interface{}{
						"banned": false,
					},
					"name": "Tom",
				},
				{
					"flags": map[string]interface{}{
						"banned": true,
					},
					"name": "Jim",
				},
			},
			[]interface{}{
				"Tom",
			},
		),
	)
}
