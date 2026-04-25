package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncToUpper is a function that converts a string to uppercase.
var FuncToUpper = NewFunc(
	"toUpper",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var err error
		if len(args) == 1 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("toUpper expects a string argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("toUpper expects data to be a string: %w", err)
			}
		}
		return model.NewStringValue(strings.ToUpper(input)), nil
	},
	ValidateArgsMax(1),
)
