package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncFlatten is a function that flattens a nested array by one level.
var FuncFlatten = NewFunc(
	"flatten",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if !input.IsSlice() {
			return nil, fmt.Errorf("flatten expects an array, got %s", input.Type().String())
		}

		res := model.NewSliceValue()
		if err := input.RangeSlice(func(i int, value *model.Value) error {
			if value.IsSlice() {
				return value.RangeSlice(func(j int, inner *model.Value) error {
					return res.Append(inner)
				})
			}
			return res.Append(value)
		}); err != nil {
			return nil, err
		}

		return res, nil
	},
	ValidateArgsMax(1),
)
