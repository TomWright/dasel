package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// deepMergeMap recursively merges next into base. When both base[key] and
// next[key] are maps the merge recurses; otherwise next[key] overwrites.
func deepMergeMap(base, next *model.Value) error {
	nextKVs, err := next.MapKeyValues()
	if err != nil {
		return err
	}
	for _, kv := range nextKVs {
		exists, _ := base.MapKeyExists(kv.Key)
		if exists && kv.Value.IsMap() {
			baseVal, _ := base.GetMapKey(kv.Key)
			if baseVal.IsMap() {
				if err := deepMergeMap(baseVal, kv.Value); err != nil {
					return err
				}
				continue
			}
		}
		if err := base.SetMapKey(kv.Key, kv.Value); err != nil {
			return err
		}
	}
	return nil
}

// FuncMerge is a function that deep-merges two or more maps together.
var FuncMerge = NewFunc(
	"merge",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		if len(args) == 1 {
			return args[0], nil
		}

		expectedType := args[0].Type()

		switch expectedType {
		case model.TypeMap:
			break
		default:
			return nil, fmt.Errorf("merge expects a map, found %s", expectedType)
		}

		// Validate types match
		for _, a := range args {
			if a.Type() != expectedType {
				return nil, fmt.Errorf("merge expects all arguments to be of the same type. expected %s, got %s", expectedType.String(), a.Type().String())
			}
		}

		base := model.NewMapValue()

		for i := 0; i < len(args); i++ {
			if err := deepMergeMap(base, args[i]); err != nil {
				return nil, fmt.Errorf("merge failed on arg %d: %w", i, err)
			}
		}

		return base, nil
	},
	ValidateArgsMin(1),
)
