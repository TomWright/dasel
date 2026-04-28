package execution

import (
	"context"
	"fmt"
	"math"

	"github.com/tomwright/dasel/v3/model"
)

// FuncCeil is a function that returns the smallest integer value greater than or equal to the input.
var FuncCeil = NewFunc(
	"ceil",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if input.IsInt() {
			v, err := input.IntValue()
			if err != nil {
				return nil, fmt.Errorf("ceil could not read int value: %w", err)
			}
			return model.NewIntValue(v), nil
		}

		v, err := input.FloatValue()
		if err != nil {
			return nil, fmt.Errorf("ceil expects a numeric value: %w", err)
		}
		return model.NewIntValue(int64(math.Ceil(v))), nil
	},
	ValidateArgsMax(1),
)
