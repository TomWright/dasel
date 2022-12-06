package dasel

import (
	"testing"
)

func TestMergeFunc(t *testing.T) {

	t.Run(
		"MergeWithArgs",
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
				[]interface{}{
					"Tom",
					"Jim",
					"Bob",
				},
			},
		),
	)

	t.Run(
		"MergeWithArgsAll",
		selectTest(
			"merge(name.first,firstNames.all()).all()",
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

	// Flaky test due to ordering.
	// t.Run(
	// 	"MergeNoArgs",
	// 	selectTest(
	// 		"name.all().merge()",
	// 		map[string]interface{}{
	// 			"name": map[string]interface{}{
	// 				"first": "Tom",
	// 				"last":  "Wright",
	// 			},
	// 		},
	// 		[]interface{}{
	// 			[]interface{}{
	// 				"Tom",
	// 				"Wright",
	// 			},
	// 		},
	// 	),
	// )

	t.Run(
		"MergeNoArgsAll",
		selectTest(
			"name.all().merge().all()",
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
}
