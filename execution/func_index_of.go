package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncIndexOf is a function that returns the index of the first occurrence of a substring.
// Returns -1 if the substring is not found.
var FuncIndexOf = NewFunc(
	"indexOf",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var substr string
		var err error

		if len(args) == 2 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("indexOf expects a string as the first argument: %w", err)
			}
			substr, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("indexOf expects a string substring as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("indexOf expects data to be a string: %w", err)
			}
			substr, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("indexOf expects a string substring as the first argument: %w", err)
			}
		}

		return model.NewIntValue(int64(strings.Index(input, substr))), nil
	},
	ValidateArgsMinMax(1, 2),
)
