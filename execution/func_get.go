package execution

import (
	"context"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
)

// FuncGet is a function returns the value at the given key/index.
var FuncGet = NewFunc(
	"get",
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
			return data.GetSliceIndex(int(index))
		case model.TypeString:
			// Given key is string, expect a map.
			if data.Type() != model.TypeMap {
				return model.NewBoolValue(false), nil
			}
			key, err := arg.StringValue()
			if err != nil {
				return nil, err
			}
			return data.GetMapKey(key)
		default:
			return nil, fmt.Errorf("get expects string or int argument")
		}
	},
	ValidateArgsMin(1),
)
