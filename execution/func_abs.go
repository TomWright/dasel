package execution

import (
	"context"
	"fmt"
	"math"

	"github.com/tomwright/dasel/v3/model"
)

// FuncAbs is a function that returns the absolute value of a number.
var FuncAbs = NewFunc(
	"abs",
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
				return nil, fmt.Errorf("abs could not read int value: %w", err)
			}
			if v < 0 {
				v = -v
			}
			return model.NewIntValue(v), nil
		}

		v, err := input.FloatValue()
		if err != nil {
			return nil, fmt.Errorf("abs expects a numeric value: %w", err)
		}
		return model.NewFloatValue(math.Abs(v)), nil
	},
	ValidateArgsMax(1),
)
