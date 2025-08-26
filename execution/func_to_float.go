package execution

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToFloat is a function that converts the given value to a string.
var FuncToFloat = NewFunc(
	"toFloat",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		switch args[0].Type() {
		case model.TypeString:
			stringValue, err := args[0].StringValue()
			if err != nil {
				return nil, err
			}

			i, err := strconv.ParseFloat(stringValue, 64)
			if err != nil {
				return nil, err
			}

			return model.NewFloatValue(i), nil
		case model.TypeInt:
			i, err := args[0].IntValue()
			if err != nil {
				return nil, err
			}
			return model.NewFloatValue(float64(i)), nil
		case model.TypeFloat:
			return args[0], nil
		case model.TypeBool:
			i, err := args[0].BoolValue()
			if err != nil {
				return nil, err
			}
			if i {
				return model.NewFloatValue(1), nil
			}
			return model.NewFloatValue(0), nil
		default:
			return nil, fmt.Errorf("cannot convert %s to float", args[0].Type())
		}
	},
	ValidateArgsExactly(1),
)
