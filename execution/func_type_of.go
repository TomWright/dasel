package execution

import (
	"github.com/tomwright/dasel/v3/model"
)

// FuncTypeOf is a function that returns the type of the first argument as a string.
var FuncTypeOf = NewFunc(
	"typeOf",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		return model.NewStringValue(args[0].Type().String()), nil
	},
	ValidateArgsExactly(1),
)
