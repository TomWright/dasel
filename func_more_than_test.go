package dasel

import (
	"testing"
)

func TestMoreThanFunc(t *testing.T) {

	t.Run("Args", selectTestErr(
		"moreThan()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "moreThan",
			Args:     []string{},
		}),
	)

	t.Run(
		"More Than",
		selectTest(
			"nums.all().moreThan(.,5)",
			map[string]interface{}{
				"nums": []any{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
				},
			},
			[]interface{}{
				false,
				false,
				false,
				false,
				false,
				false,
				true,
				true,
				true,
				true,
			},
		),
	)
}
