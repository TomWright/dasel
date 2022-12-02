package dasel

import (
	"testing"
)

func TestLessThanFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"lessThan()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "lessThan",
			Args:     []string{},
		}),
	)

	t.Run(
		"Less Than",
		selectTest(
			"nums.all().lessThan(.,5)",
			map[string]interface{}{
				"nums": []any{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
				},
			},
			[]interface{}{
				true,
				true,
				true,
				true,
				true,
				false,
				false,
				false,
				false,
				false,
			},
		),
	)
}
