package execution

import (
	"context"
	"fmt"
	"math"

	"github.com/tomwright/dasel/v3/model"
)

// FuncFloor is a function that returns the largest integer value less than or equal to the input.
var FuncFloor = NewFunc(
	"floor",
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
				return nil, fmt.Errorf("floor could not read int value: %w", err)
			}
			return model.NewIntValue(v), nil
		}

		v, err := input.FloatValue()
		if err != nil {
			return nil, fmt.Errorf("floor expects a numeric value: %w", err)
		}
		return model.NewIntValue(int64(math.Floor(v))), nil
	},
	ValidateArgsMax(1),
)
