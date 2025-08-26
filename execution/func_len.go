package execution

import (
	"context"
	"github.com/tomwright/dasel/v3/model"
)

// FuncLen is a function that returns the length of the given value.
var FuncLen = NewFunc(
	"len",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		arg := args[0]

		l, err := arg.Len()
		if err != nil {
			return nil, err
		}

		return model.NewIntValue(int64(l)), nil
	},
	ValidateArgsExactly(1),
)
