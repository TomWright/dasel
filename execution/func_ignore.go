package execution

import (
	"context"
	"github.com/tomwright/dasel/v3/model"
)

// FuncIgnore is a function that ignores the value, causing it to be rejected from a branch.
var FuncIgnore = NewFunc(
	"ignore",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		data.MarkAsIgnore()
		return data, nil
	},
	ValidateArgsExactly(0),
)
