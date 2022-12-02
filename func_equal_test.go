package dasel

import (
	"testing"
)

func TestEqualFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"equal()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "equal",
			Args:     []string{},
		}),
	)

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
				true,
				false,
			},
		),
	)

	t.Run(
		"Multi Equal",
		selectTest(
			"name.all().equal(key(),first,key(),first)",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{
				true,
				false,
			},
		),
	)

	t.Run(
		"Single Equal Optional Field",
		selectTest(
			"all().equal(primary,true)",
			[]interface{}{
				map[string]interface{}{
					"name":    "red",
					"hex":     "ff0000",
					"primary": true,
				},
				map[string]interface{}{
					"name":    "green",
					"hex":     "00ff00",
					"primary": true,
				},
				map[string]interface{}{
					"name":    "blue",
					"hex":     "0000ff",
					"primary": true,
				},
				map[string]interface{}{
					"name":    "orange",
					"hex":     "ffa500",
					"primary": false,
				},
			},
			[]interface{}{
				true, true, true, false,
			},
		),
	)
}
