package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncAvg is a function that returns the average of the given numbers.
// Always returns a float value.
var FuncAvg = NewFunc(
	"avg",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var sum float64

		for _, arg := range args {
			if arg.IsInt() {
				intVal, err := arg.IntValue()
				if err != nil {
					return nil, err
				}
				sum += float64(intVal)
				continue
			}
			if arg.IsFloat() {
				floatVal, err := arg.FloatValue()
				if err != nil {
					return nil, err
				}
				sum += floatVal
				continue
			}
			return nil, fmt.Errorf("cannot average non-numeric value of type %s", arg.Type().String())
		}

		return model.NewFloatValue(sum / float64(len(args))), nil
	},
	ValidateArgsMin(1),
)
