package execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncJoin is a function that returns the joins the given data or args to a string.
var FuncJoin = NewFunc(
	"join",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		separator, err := args[0].StringValue()
		if err != nil {
			return nil, fmt.Errorf("join expects a string separator as the first argument: %w", err)
		}

		var valuesToJoin []string

		if len(args) == 2 && args[1].IsSlice() {
			if err := args[1].RangeSlice(func(i int, value *model.Value) error {
				strVal, err := value.StringValue()
				if err != nil {
					return fmt.Errorf("could not read string value of index %d: %w", i, err)
				}
				valuesToJoin = append(valuesToJoin, strVal)
				return nil
			}); err != nil {
				return nil, err
			}
		} else if len(args) > 1 {
			// Join the args
			for i := 1; i < len(args); i++ {
				strVal, err := args[i].StringValue()
				if err != nil {
					return nil, fmt.Errorf("could not read string value of argument index %d: %w", i, err)
				}
				valuesToJoin = append(valuesToJoin, strVal)
			}
		} else {
			if err := data.RangeSlice(func(i int, value *model.Value) error {
				strVal, err := value.StringValue()
				if err != nil {
					return fmt.Errorf("could not read string value of index %d: %w", i, err)
				}
				valuesToJoin = append(valuesToJoin, strVal)
				return nil
			}); err != nil {
				return nil, err
			}
		}

		joined := strings.Join(valuesToJoin, separator)

		return model.NewStringValue(joined), nil
	},
	ValidateArgsMin(1),
)
