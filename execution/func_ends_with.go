package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncEndsWith is a function that checks if a string ends with a given suffix.
var FuncEndsWith = NewFunc(
	"endsWith",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var suffix string
		var err error

		if len(args) == 2 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("endsWith expects a string as the first argument: %w", err)
			}
			suffix, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("endsWith expects a string suffix as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("endsWith expects data to be a string: %w", err)
			}
			suffix, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("endsWith expects a string suffix as the first argument: %w", err)
			}
		}

		return model.NewBoolValue(strings.HasSuffix(input, suffix)), nil
	},
	ValidateArgsMinMax(1, 2),
)
