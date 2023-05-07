package dasel

import (
	"github.com/tomwright/dasel/v2/dencoding"
	"strings"
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

	original := dencoding.NewMap().
		Set("name", dencoding.NewMap().
			Set("first", "Tom").
			Set("last", "Wright")).
		Set("colours", []interface{}{
			"red", "green", "blue",
		})

	t.Run(
		"JoinCommaSeparator",
		selectTestAssert(
			"name.all().join(\\,)",
			original,
			func(t *testing.T, got []any) {
				required := []string{"Tom", "Wright"}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
				str, ok := got[0].(string)
				if !ok {
					t.Errorf("expected 1st result to be a string, got %T", got[0])
					return
				}

				gotStrs := strings.Split(str, ",")
				for _, req := range required {
					found := false
					for _, got := range gotStrs {
						if got == req {
							found = true
							continue
						}
					}
					if !found {
						t.Errorf("expected %v, got %v", required, got)
					}
				}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
			},
		),
	)

	t.Run(
		"JoinSpaceSeparator",
		selectTestAssert(
			"name.all().join( )",
			original,
			func(t *testing.T, got []any) {
				required := []string{"Tom", "Wright"}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
				str, ok := got[0].(string)
				if !ok {
					t.Errorf("expected 1st result to be a string, got %T", got[0])
					return
				}

				gotStrs := strings.Split(str, " ")
				for _, req := range required {
					found := false
					for _, got := range gotStrs {
						if got == req {
							found = true
							continue
						}
					}
					if !found {
						t.Errorf("expected %v, got %v", required, got)
					}
				}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
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
		selectTestAssert(
			"all().join(\\,,all())",
			dencoding.NewMap().
				Set("x", []interface{}{1, 2, 3}).
				Set("y", []interface{}{4, 5, 6}).
				Set("z", []interface{}{7, 8, 9}),
			func(t *testing.T, got []any) {
				required := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
				str, ok := got[0].(string)
				if !ok {
					t.Errorf("expected 1st result to be a string, got %T", got[0])
					return
				}

				gotStrs := strings.Split(str, ",")
				for _, req := range required {
					found := false
					for _, got := range gotStrs {
						if got == req {
							found = true
							continue
						}
					}
					if !found {
						t.Errorf("expected %v, got %v", required, got)
					}
				}
				if len(got) != 1 {
					t.Errorf("expected 1 result, got %v", got)
					return
				}
			},
		),
	)
}
