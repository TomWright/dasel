package dasel

import (
	"testing"
)

func TestMergeFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"merge()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "merge",
			Args:     []string{},
		}),
	)

	t.Run(
		"Merge",
		selectTest(
			"merge(name.first,firstNames.all())",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
				"firstNames": []interface{}{
					"Jim",
					"Bob",
				},
			},
			[]interface{}{
				"Tom",
				"Jim",
				"Bob",
			},
		),
	)
}
