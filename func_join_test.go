package dasel

import (
	"github.com/tomwright/dasel/v2/ordered"
	"testing"
)

func TestJoinFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"join()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "join",
			Args:     []string{},
		}),
	)

	original := ordered.NewMap().
		Set("name", ordered.NewMap().
			Set("first", "Tom").
			Set("last", "Wright")).
		Set("colours", []interface{}{
			"red", "green", "blue",
		})

	t.Run(
		"JoinCommaSeparator",
		selectTest(
			"name.all().join(\\,)",
			original,
			[]interface{}{
				"Tom,Wright",
			},
		),
	)

	t.Run(
		"JoinSpaceSeparator",
		selectTest(
			"name.all().join( )",
			original,
			[]interface{}{
				"Tom Wright",
			},
		),
	)

	t.Run(
		"JoinWithSeparatorsAndSelectors",
		selectTest(
			"name.join( ,last,first)",
			original,
			[]interface{}{
				"Wright Tom",
			},
		),
	)

	t.Run(
		"JoinInMap",
		selectTest(
			"mapOf(first,name.first,last,name.last,full,name.join( ,string(Mr),first,last))",
			original,
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"full":  "Mr Tom Wright",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"JoinManyLists",
		selectTest(
			"all().join(\\,,all())",
			ordered.NewMap().
				Set("x", []interface{}{1, 2, 3}).
				Set("y", []interface{}{4, 5, 6}).
				Set("z", []interface{}{7, 8, 9}),
			[]interface{}{
				"1,2,3,4,5,6,7,8,9",
			},
		),
	)
}
