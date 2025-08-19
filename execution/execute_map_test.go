package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestMap(t *testing.T) {
	t.Run("property from slice of maps", testCase{
		inFn: func() *model.Value {
			return model.NewValue([]any{
				orderedmap.NewMap().Set("number", 1),
				orderedmap.NewMap().Set("number", 2),
				orderedmap.NewMap().Set("number", 3),
			})
		},
		s: `map(number)`,
		outFn: func() *model.Value {
			return model.NewValue([]any{1, 2, 3})
		},
	}.run)
	t.Run("with chain of selectors", testCase{
		inFn: func() *model.Value {
			return model.NewValue([]any{
				orderedmap.NewMap().Set("foo", 1).Set("bar", 4),
				orderedmap.NewMap().Set("foo", 2).Set("bar", 5),
				orderedmap.NewMap().Set("foo", 3).Set("bar", 6),
			})
		},
		s: `
				map (
					{
						total: add( foo, bar, 1 )
					}
				)
				.map ( total )`,
		outFn: func() *model.Value {
			res := model.NewSliceValue()
			if err := res.Append(model.NewValue(6)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewValue(8)); err != nil {
				t.Fatal(err)
			}
			if err := res.Append(model.NewValue(10)); err != nil {
				t.Fatal(err)
			}
			return res
		},
	}.run)
}
