package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncUnique is a function that removes duplicate values from an array.
var FuncUnique = NewFunc(
	"unique",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if !input.IsSlice() {
			return nil, fmt.Errorf("unique expects an array, got %s", input.Type().String())
		}

		res := model.NewSliceValue()
		if err := input.RangeSlice(func(i int, value *model.Value) error {
			duplicate := false
			_ = res.RangeSlice(func(j int, existing *model.Value) error {
				equal, err := value.EqualTypeValue(existing)
				if err != nil {
					return nil
				}
				if equal {
					duplicate = true
				}
				return nil
			})
			if !duplicate {
				return res.Append(value)
			}
			return nil
		}); err != nil {
			return nil, err
		}

		return res, nil
	},
	ValidateArgsMax(1),
)
