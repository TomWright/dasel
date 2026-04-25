package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToBool is a function that converts the given value to a bool.
var FuncToBool = NewFunc(
	"toBool",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		switch args[0].Type() {
		case model.TypeBool:
			b, err := args[0].BoolValue()
			if err != nil {
				return nil, err
			}
			return model.NewBoolValue(b), nil
		case model.TypeString:
			s, err := args[0].StringValue()
			if err != nil {
				return nil, err
			}
			switch strings.ToLower(s) {
			case "true", "1", "yes":
				return model.NewBoolValue(true), nil
			case "false", "0", "no", "":
				return model.NewBoolValue(false), nil
			default:
				return nil, fmt.Errorf("cannot convert string %q to bool", s)
			}
		case model.TypeInt:
			i, err := args[0].IntValue()
			if err != nil {
				return nil, err
			}
			return model.NewBoolValue(i != 0), nil
		case model.TypeFloat:
			f, err := args[0].FloatValue()
			if err != nil {
				return nil, err
			}
			return model.NewBoolValue(f != 0), nil
		case model.TypeNull:
			return model.NewBoolValue(false), nil
		default:
			return nil, fmt.Errorf("cannot convert %s to bool", args[0].Type())
		}
	},
	ValidateArgsExactly(1),
)
