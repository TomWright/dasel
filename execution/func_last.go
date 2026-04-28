package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncLast is a function that returns the last element of an array.
var FuncLast = NewFunc(
	"last",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if !input.IsSlice() {
			return nil, fmt.Errorf("last expects an array, got %s", input.Type().String())
		}

		length, err := input.SliceLen()
		if err != nil {
			return nil, err
		}
		if length == 0 {
			return model.NewNullValue(), nil
		}

		return input.GetSliceIndex(length - 1)
	},
	ValidateArgsMax(1),
)
