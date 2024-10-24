package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncMerge is a function that merges two or more items together.
var FuncMerge = NewFunc(
	"merge",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		if len(args) == 1 {
			return args[0], nil
		}

		expectedType := args[0].Type()

		switch expectedType {
		case model.TypeMap:
			break
		default:
			return nil, fmt.Errorf("merge exects a map, found %s", expectedType)
		}

		// Validate types match
		for _, a := range args {
			if a.Type() != expectedType {
				return nil, fmt.Errorf("merge expects all arguments to be of the same type. expected %s, got %s", expectedType.String(), a.Type().String())
			}
		}

		base := model.NewMapValue()

		for i := 0; i < len(args); i++ {
			next := args[i]

			nextKVs, err := next.MapKeyValues()
			if err != nil {
				return nil, fmt.Errorf("merge failed to extract key values for arg %d: %w", i, err)
			}

			for _, kv := range nextKVs {
				if err := base.SetMapKey(kv.Key, kv.Value); err != nil {
					return nil, fmt.Errorf("merge failed to set map key %s: %w", kv.Key, err)
				}
			}
		}

		return base, nil
	},
	ValidateArgsMin(1),
)
