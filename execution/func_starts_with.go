package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncStartsWith is a function that checks if a string starts with a given prefix.
var FuncStartsWith = NewFunc(
	"startsWith",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var prefix string
		var err error

		if len(args) == 2 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("startsWith expects a string as the first argument: %w", err)
			}
			prefix, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("startsWith expects a string prefix as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("startsWith expects data to be a string: %w", err)
			}
			prefix, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("startsWith expects a string prefix as the first argument: %w", err)
			}
		}

		return model.NewBoolValue(strings.HasPrefix(input, prefix)), nil
	},
	ValidateArgsMinMax(1, 2),
)
