package dasel

import (
	"testing"
)

func TestTypeFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"type(x)",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "type",
			Args:     []string{"x"},
		}),
	)

	data := map[string]interface{}{
		"string": "hello",
		"slice": []interface{}{
			1, 2, 3,
		},
		"map": map[string]interface{}{
			"x": 1,
		},
		"int":   int(1),
		"float": float32(1),
		"bool":  true,
	}

	t.Run(
		"String",
		selectTest(
			"string.type()",
			data,
			[]interface{}{
				"string",
			},
		),
	)
	t.Run(
		"Slice",
		selectTest(
			"slice.type()",
			data,
			[]interface{}{
				"array",
			},
		),
	)
	t.Run(
		"map",
		selectTest(
			"map.type()",
			data,
			[]interface{}{
				"object",
			},
		),
	)
	t.Run(
		"int",
		selectTest(
			"int.type()",
			data,
			[]interface{}{
				"number",
			},
		),
	)
	t.Run(
		"float",
		selectTest(
			"float.type()",
			data,
			[]interface{}{
				"number",
			},
		),
	)
	t.Run(
		"bool",
		selectTest(
			"bool.type()",
			data,
			[]interface{}{
				"bool",
			},
		),
	)
}
