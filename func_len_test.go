package dasel

import (
	"testing"
)

func TestLenFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"len(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "len",
			Args:     []string{"x"},
		}),
	)

	data := map[string]interface{}{
		"string": "hello",
		"slice": []interface{}{
			1, 2, 3,
		},
		"falseBool": false,
		"trueBool":  true,
	}

	t.Run(
		"String",
		selectTest(
			"string.len()",
			data,
			[]interface{}{5},
		),
	)
	t.Run(
		"Slice",
		selectTest(
			"slice.len()",
			data,
			[]interface{}{3},
		),
	)
	t.Run(
		"False Bool",
		selectTest(
			"falseBool.len()",
			data,
			[]interface{}{0},
		),
	)
	t.Run(
		"True Bool",
		selectTest(
			"trueBool.len()",
			data,
			[]interface{}{1},
		),
	)
}
