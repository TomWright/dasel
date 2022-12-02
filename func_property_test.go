package dasel

import (
	"testing"
)

func TestPropertyFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"property()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "property",
			Args:     []string{},
		}),
	)

	t.Run("NotFound", selectTestErr(
		"asd",
		map[string]interface{}{"x": "y"},
		&ErrPropertyNotFound{
			Property: "asd",
		}),
	)

	t.Run("NotFoundOnString", selectTestErr(
		"x.asd",
		map[string]interface{}{"x": "y"},
		&ErrPropertyNotFound{
			Property: "asd",
		}),
	)

	original := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tom",
			"last":  "Wright",
		},
		"colours": []interface{}{
			"red", "green", "blue",
		},
	}

	t.Run(
		"SingleLevelProperty",
		selectTest(
			"name",
			original,
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"SingleLevelPropertyFunc",
		selectTest(
			"property(name)",
			original,
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"NestedPropertyFunc",
		selectTest(
			"property(name).property(first)",
			original,
			[]interface{}{
				"Tom",
			},
		),
	)

	t.Run(
		"NestedMultiPropertyFunc",
		selectTest(
			"property(name).property(first,last)",
			original,
			[]interface{}{
				"Tom",
				"Wright",
			},
		),
	)

	t.Run(
		"NestedMultiMissingPropertyFunc",
		selectTest(
			"property(name).property(first,last,middle?)",
			original,
			[]interface{}{
				"Tom",
				"Wright",
			},
		),
	)
}
