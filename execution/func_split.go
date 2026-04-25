package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncSplit is a function that splits a string by a separator into an array.
var FuncSplit = NewFunc(
	"split",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		separator, err := args[0].StringValue()
		if err != nil {
			return nil, fmt.Errorf("split expects a string separator as the first argument: %w", err)
		}

		var input string
		if len(args) == 2 {
			input, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("split expects a string as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("split expects data to be a string: %w", err)
			}
		}

		parts := strings.Split(input, separator)

		res := model.NewSliceValue()
		for _, part := range parts {
			if err := res.Append(model.NewStringValue(part)); err != nil {
				return nil, fmt.Errorf("could not append split result: %w", err)
			}
		}

		return res, nil
	},
	ValidateArgsMinMax(1, 2),
)
