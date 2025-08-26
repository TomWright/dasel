package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToString is a function that converts the given value to a string.
var FuncToString = NewFunc(
	"toString",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		switch args[0].Type() {
		case model.TypeString:
			stringValue, err := args[0].StringValue()
			if err != nil {
				return nil, err
			}
			model.NewStringValue(stringValue)
			return args[0], nil
		case model.TypeInt:
			i, err := args[0].IntValue()
			if err != nil {
				return nil, err
			}
			return model.NewStringValue(fmt.Sprintf("%d", i)), nil
		case model.TypeFloat:
			i, err := args[0].FloatValue()
			if err != nil {
				return nil, err
			}
			return model.NewStringValue(fmt.Sprintf("%g", i)), nil
		case model.TypeBool:
			i, err := args[0].BoolValue()
			if err != nil {
				return nil, err
			}
			return model.NewStringValue(fmt.Sprintf("%v", i)), nil
		default:
			return nil, fmt.Errorf("cannot convert %s to string", args[0].Type())
		}
	},
	ValidateArgsExactly(1),
)
