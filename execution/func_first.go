package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncFirst is a function that returns the first element of an array.
var FuncFirst = NewFunc(
	"first",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if !input.IsSlice() {
			return nil, fmt.Errorf("first expects an array, got %s", input.Type().String())
		}

		length, err := input.SliceLen()
		if err != nil {
			return nil, err
		}
		if length == 0 {
			return model.NewNullValue(), nil
		}

		return input.GetSliceIndex(0)
	},
	ValidateArgsMax(1),
)
