package dasel

import (
	"testing"
)

func TestFilterFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"filter()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "filter",
			Args:     []string{},
		}),
	)

	t.Run(
		"Filter Equal Key",
		selectTest(
			"name.all().filter(equal(key(),first))",
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
		"Multiple Filter Equal Key",
		selectTest(
			"name.all().filter(equal(key(),first),equal(key(),last))",
			map[string]interface{}{
				"name": map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
			[]interface{}{},
		),
	)

	t.Run(
		"Filter Equal Prop",
		selectTest(
			"all().filter(equal(primary,true)).name",
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
				"red", "green", "blue",
			},
		),
	)

	t.Run(
		"FilterNestedProp",
		selectTest(
			"all().filter(equal(flags.banned,false)).name",
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

	t.Run(
		"Filter And",
		selectTest(
			"all().filter(and(equal(primary,true),equal(name,red))).name",
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
				"red",
			},
		),
	)

	t.Run(
		"Filter And",
		selectTest(
			"all().filter(and(equal(primary,true),equal(name,orange))).name",
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
			[]interface{}{},
		),
	)

	t.Run(
		"Filter Or",
		selectTest(
			"all().filter(or(equal(primary,true),equal(name,orange))).name",
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
				"red", "green", "blue", "orange",
			},
		),
	)
}
