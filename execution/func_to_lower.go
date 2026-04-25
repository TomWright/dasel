package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToLower is a function that converts a string to lowercase.
var FuncToLower = NewFunc(
	"toLower",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var err error
		if len(args) == 1 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("toLower expects a string argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("toLower expects data to be a string: %w", err)
			}
		}
		return model.NewStringValue(strings.ToLower(input)), nil
	},
	ValidateArgsMax(1),
)
