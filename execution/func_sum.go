package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncSum is a function that returns the sum of the given numbers.
var FuncSum = NewFunc(
	"sum",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		returnType := model.TypeInt

		for _, arg := range args {
			if arg.IsInt() {
				continue
			}
			if arg.IsFloat() {
				returnType = model.TypeFloat
				break
			}
			return nil, fmt.Errorf("cannot sum non-numeric value of type %s", arg.Type().String())
		}

		switch returnType {
		case model.TypeInt:
			var sum int64
			for _, arg := range args {
				if arg.IsInt() {
					intVal, err := arg.IntValue()
					if err != nil {
						return nil, err
					}
					sum += intVal
					continue
				}

				floatVal, err := arg.FloatValue()
				if err != nil {
					return nil, err
				}
				sum += int64(floatVal)
			}
			return model.NewIntValue(sum), nil
		case model.TypeFloat:
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

				floatVal, err := arg.FloatValue()
				if err != nil {
					return nil, err
				}
				sum += floatVal
			}
			return model.NewFloatValue(sum), nil
		default:
			return nil, fmt.Errorf("unsupported return type %s", returnType.String())
		}
	},
	ValidateArgsMin(1),
)
