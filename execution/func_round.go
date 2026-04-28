package execution

import (
	"context"
	"fmt"
	"math"

	"github.com/tomwright/dasel/v3/model"
)

// FuncRound is a function that rounds a number to the nearest integer.
var FuncRound = NewFunc(
	"round",
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
				return nil, fmt.Errorf("round could not read int value: %w", err)
			}
			return model.NewIntValue(v), nil
		}

		v, err := input.FloatValue()
		if err != nil {
			return nil, fmt.Errorf("round expects a numeric value: %w", err)
		}
		return model.NewIntValue(int64(math.Round(v))), nil
	},
	ValidateArgsMax(1),
)
