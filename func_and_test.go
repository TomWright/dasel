package dasel

import (
	"testing"
)

func TestAndFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"and()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "and",
			Args:     []string{},
		}),
	)

	t.Run(
		"NoneEqualMoreThan",
		selectTest(
			"numbers.all().and(equal(.,2),moreThan(.,2))",
			map[string]interface{}{
				"numbers": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			[]interface{}{
				false, false, false, false, false, false, false, false, false, false,
			},
		),
	)
	t.Run(
		"SomeEqualMoreThan",
		selectTest(
			"numbers.all().and(equal(.,4),moreThan(.,2))",
			map[string]interface{}{
				"numbers": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			[]interface{}{
				false, false, false, false, true, false, false, false, false, false,
			},
		),
	)
}
