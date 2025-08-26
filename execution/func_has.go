package execution

import (
	"context"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
)

// FuncHas is a function that true or false if the input has the given key/index.
var FuncHas = NewFunc(
	"has",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {

		arg := args[0]

		switch arg.Type() {
		case model.TypeInt:
			// Given key is int, expect a slice.
			if data.Type() != model.TypeSlice {
				return model.NewBoolValue(false), nil
			}
			index, err := arg.IntValue()
			if err != nil {
				return nil, err
			}
			sliceLen, err := data.SliceLen()
			if err != nil {
				return nil, err
			}
			return model.NewBoolValue(index >= 0 && index < int64(sliceLen)), nil
		case model.TypeString:
			// Given key is string, expect a map.
			if data.Type() != model.TypeMap {
				return model.NewBoolValue(false), nil
			}
			key, err := arg.StringValue()
			if err != nil {
				return nil, err
			}
			exists, err := data.MapKeyExists(key)
			if err != nil {
				return nil, err
			}
			return model.NewBoolValue(exists), nil
		default:
			return nil, fmt.Errorf("has expects string or int argument")
		}
	},
	ValidateArgsMin(1),
)
