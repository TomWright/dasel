package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncEntries converts a map into an array of {key, value} objects.
var FuncEntries = NewFunc(
	"entries",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		if !data.IsMap() {
			return nil, fmt.Errorf("entries can only be used on maps, got %s", data.Type().String())
		}

		res := model.NewSliceValue()
		if err := data.RangeMap(func(key string, value *model.Value) error {
			entry := model.NewMapValue()
			if err := entry.SetMapKey("key", model.NewStringValue(key)); err != nil {
				return err
			}
			if err := entry.SetMapKey("value", value); err != nil {
				return err
			}
			return res.Append(entry)
		}); err != nil {
			return nil, err
		}

		return res, nil
	},
	ValidateArgsExactly(0),
)

// FuncFromEntries converts an array of {key, value} objects into a map.
var FuncFromEntries = NewFunc(
	"fromEntries",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var input *model.Value
		if len(args) == 1 {
			input = args[0]
		} else {
			input = data
		}

		if !input.IsSlice() {
			return nil, fmt.Errorf("fromEntries expects an array, got %s", input.Type().String())
		}

		res := model.NewMapValue()
		if err := input.RangeSlice(func(i int, entry *model.Value) error {
			if !entry.IsMap() {
				return fmt.Errorf("fromEntries expects each element to be an object, got %s at index %d", entry.Type().String(), i)
			}

			keyVal, err := entry.GetMapKey("key")
			if err != nil {
				return fmt.Errorf("fromEntries expects each element to have a \"key\" field: %w", err)
			}
			key, err := keyVal.StringValue()
			if err != nil {
				return fmt.Errorf("fromEntries expects \"key\" to be a string: %w", err)
			}

			value, err := entry.GetMapKey("value")
			if err != nil {
				return fmt.Errorf("fromEntries expects each element to have a \"value\" field: %w", err)
			}

			return res.SetMapKey(key, value)
		}); err != nil {
			return nil, err
		}

		return res, nil
	},
	ValidateArgsMax(1),
)
