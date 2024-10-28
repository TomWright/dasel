package execution

import (
	"fmt"
	"strconv"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToInt is a function that converts the given value to a string.
var FuncToInt = NewFunc(
	"toInt",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		switch args[0].Type() {
		case model.TypeString:
			stringValue, err := args[0].StringValue()
			if err != nil {
				return nil, err
			}

			i, err := strconv.ParseInt(stringValue, 10, 64)
			if err != nil {
				return nil, err
			}

			return model.NewIntValue(i), nil
		case model.TypeInt:
			return args[0], nil
		case model.TypeFloat:
			i, err := args[0].FloatValue()
			if err != nil {
				return nil, err
			}
			return model.NewIntValue(int64(i)), nil
		case model.TypeBool:
			i, err := args[0].BoolValue()
			if err != nil {
				return nil, err
			}
			if i {
				return model.NewIntValue(1), nil
			}
			return model.NewIntValue(0), nil
		default:
			return nil, fmt.Errorf("cannot convert %s to int", args[0].Type())
		}
	},
	ValidateArgsExactly(1),
)
