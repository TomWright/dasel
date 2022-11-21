package dasel

import (
	"testing"
)

func TestCountFunc(t *testing.T) {
	data := map[string]interface{}{
		"string": "hello",
		"slice": []interface{}{
			1, 2, 3,
		},
		"falseBool": false,
		"trueBool":  true,
	}

	t.Run(
		"RootObject",
		selectTest(
			"count()",
			data,
			[]interface{}{1},
		),
	)
	t.Run(
		"All",
		selectTest(
			"all().count()",
			data,
			[]interface{}{4},
		),
	)
	t.Run(
		"NestedAll",
		selectTest(
			"slice.all().count()",
			data,
			[]interface{}{3},
		),
	)
}
