package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncReverse is a function that reverses the input.
var FuncReverse = NewFunc(
	"reverse",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		arg := args[0]

		switch arg.Type() {
		case model.TypeString:
			return arg.StringIndexRange(-1, 0)
		case model.TypeSlice:
			return arg.SliceIndexRange(-1, 0)
		default:
			return nil, fmt.Errorf("reverse expects a slice or string, got %s", arg.Type())
		}
	},
	ValidateArgsExactly(1),
)
