package dasel

import (
	"testing"
)

func TestOrDefaultFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"orDefault()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "property",
			Args:     []string{},
		}),
	)

	t.Run("OriginalAndDefaultNotFoundProperty", selectTestErr(
		"orDefault(a,b)",
		map[string]interface{}{"x": "y"},
		&ErrPropertyNotFound{
			Property: "b",
		}),
	)

	t.Run("OriginalAndDefaultNotFoundIndex", selectTestErr(
		"orDefault(x.[1],x.[2])",
		map[string]interface{}{"x": []int{1}},
		&ErrIndexNotFound{
			Index: 2,
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
		"FirstNameOrLastName",
		selectTest(
			"orDefault(name.first,name.last)",
			original,
			[]interface{}{
				"Tom",
			},
		),
	)

	t.Run(
		"MiddleNameOrDefault",
		selectTest(
			"orDefault(name.middle,string(default))",
			original,
			[]interface{}{
				"default",
			},
		),
	)

	t.Run(
		"FirstColourOrSecondColour",
		selectTest(
			"orDefault(colours.[0],colours.[2])",
			original,
			[]interface{}{
				"red",
			},
		),
	)

	t.Run(
		"FourthColourOrDefault",
		selectTest(
			"orDefault(colours.[3],string(default))",
			original,
			[]interface{}{
				"default",
			},
		),
	)
}
