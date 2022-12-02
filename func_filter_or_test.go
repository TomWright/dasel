package dasel

import (
	"testing"
)

func TestFilterOrFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"filterOr()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "filterOr",
			Args:     []string{},
		}),
	)

	t.Run(
		"Filter Equal Key",
		selectTest(
			"name.all().filterOr(equal(key(),first))",
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
		"Multiple Filter Or Equal Key",
		selectTest(
			"name.all().filterOr(equal(key(),first),equal(key(),last))",
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
		"MoreThanEqual",
		selectTest(
			"nums.all().filterOr(moreThan(.,3),equal(.,3))",
			map[string]interface{}{
				"nums": []interface{}{0, 1, 2, 3, 4, 5},
			},
			[]interface{}{3, 4, 5},
		),
	)

	t.Run(
		"LessThanEqual",
		selectTest(
			"nums.all().filterOr(lessThan(.,3),equal(.,3))",
			map[string]interface{}{
				"nums": []interface{}{0, 1, 2, 3, 4, 5},
			},
			[]interface{}{0, 1, 2, 3},
		),
	)
}
