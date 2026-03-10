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
			index, err := arg.IntValue()
			if err != nil {
				return nil, err
			}
			return data.GetSliceIndex(int(index))
		case model.TypeString:
			key, err := arg.StringValue()
			if err != nil {
				return nil, err
			}
			return data.GetMapKey(key)
		default:
			return nil, fmt.Errorf("get expects string or int argument")
		}
	},
	ValidateArgsExactly(1),
)
