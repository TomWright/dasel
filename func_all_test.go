package dasel

import "testing"

func TestAllFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"all(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "all",
			Args:     []string{"x"},
		}),
	)

	t.Run(
		"RootAll",
		selectTest(
			"all()",
			[]interface{}{"red", "green", "blue"},
			[]interface{}{"red", "green", "blue"},
		),
	)
	t.Run(
		"NestedAll",
		selectTest(
			"colours.all()",
			map[string]interface{}{
				"colours": []interface{}{"red", "green", "blue"},
			},
			[]interface{}{"red", "green", "blue"},
		),
	)
}
