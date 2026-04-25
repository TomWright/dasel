package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncTrim is a function that trims whitespace from both ends of a string.
var FuncTrim = NewFunc(
	"trim",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var err error
		if len(args) == 1 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trim expects a string argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("trim expects data to be a string: %w", err)
			}
		}
		return model.NewStringValue(strings.TrimSpace(input)), nil
	},
	ValidateArgsMax(1),
)

// FuncTrimPrefix is a function that trims a prefix from a string.
var FuncTrimPrefix = NewFunc(
	"trimPrefix",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var prefix string
		var err error

		if len(args) == 2 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimPrefix expects a string as the first argument: %w", err)
			}
			prefix, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimPrefix expects a string prefix as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimPrefix expects data to be a string: %w", err)
			}
			prefix, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimPrefix expects a string prefix as the first argument: %w", err)
			}
		}

		return model.NewStringValue(strings.TrimPrefix(input, prefix)), nil
	},
	ValidateArgsMinMax(1, 2),
)

// FuncTrimSuffix is a function that trims a suffix from a string.
var FuncTrimSuffix = NewFunc(
	"trimSuffix",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input string
		var suffix string
		var err error

		if len(args) == 2 {
			input, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimSuffix expects a string as the first argument: %w", err)
			}
			suffix, err = args[1].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimSuffix expects a string suffix as the second argument: %w", err)
			}
		} else {
			input, err = data.StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimSuffix expects data to be a string: %w", err)
			}
			suffix, err = args[0].StringValue()
			if err != nil {
				return nil, fmt.Errorf("trimSuffix expects a string suffix as the first argument: %w", err)
			}
		}

		return model.NewStringValue(strings.TrimSuffix(input, suffix)), nil
	},
	ValidateArgsMinMax(1, 2),
)
