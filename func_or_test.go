package dasel

import (
	"testing"
)

func TestOrFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"or()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "or",
			Args:     []string{},
		}),
	)

	t.Run(
		"NoneEqualMoreThan",
		selectTest(
			"numbers.all().or(equal(.,2),moreThan(.,2))",
			map[string]interface{}{
				"numbers": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			[]interface{}{
				false, false, true, true, true, true, true, true, true, true,
			},
		),
	)
	t.Run(
		"SomeEqualMoreThan",
		selectTest(
			"numbers.all().or(equal(.,0),moreThan(.,2))",
			map[string]interface{}{
				"numbers": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			[]interface{}{
				true, false, false, true, true, true, true, true, true, true,
			},
		),
	)
}
