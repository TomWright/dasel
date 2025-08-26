package execution

import (
	"context"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
)

// FuncContains is a function that returns the highest number.
var FuncContains = NewFunc(
	"contains",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var contains bool

		target := args[0]

		length, err := data.SliceLen()
		if err != nil {
			return nil, fmt.Errorf("error getting slice length: %w", err)
		}

		for i := 0; i < length; i++ {
			v, err := data.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index %d: %w", i, err)
			}
			matches, err := v.Equal(target)
			if err != nil {
				continue
			}
			matchesBool, err := matches.BoolValue()
			if err != nil {
				return nil, err
			}
			if matchesBool {
				contains = true
				break
			}
		}

		return model.NewBoolValue(contains), nil
	},
	ValidateArgsExactly(1),
)
