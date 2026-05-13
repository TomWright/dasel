package kdl

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/kdl/internal"
)

var _ parsing.Reader = (*kdlReader)(nil)

func newKDLReader(_ parsing.ReaderOptions) (parsing.Reader, error) {
	return &kdlReader{}, nil
}

type kdlReader struct{}

func (r *kdlReader) Read(data []byte) (*model.Value, error) {
	doc, err := internal.Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("kdl: parse error: %w", err)
	}
	return nodesToValue(doc.Nodes)
}

// nodesToValue converts a list of KDL nodes into a model.Value map.
// Duplicate node names at the same level are promoted to slices.
func nodesToValue(nodes []*internal.Node) (*model.Value, error) {
	result := model.NewMapValue()
	seen := make(map[string]int) // tracks how many times a name appears

	for _, node := range nodes {
		val, err := nodeToValue(node)
		if err != nil {
			return nil, err
		}

		count := seen[node.Name]
		seen[node.Name] = count + 1

		switch count {
		case 0:
			// First time seeing this name
			if err := result.SetMapKey(node.Name, val); err != nil {
				return nil, err
			}
		case 1:
			// Second time — promote existing value to slice
			existing, err := result.GetMapKey(node.Name)
			if err != nil {
				return nil, err
			}
			slice := model.NewSliceValue()
			if err := slice.Append(existing); err != nil {
				return nil, err
			}
			if err := slice.Append(val); err != nil {
				return nil, err
			}
			if err := result.SetMapKey(node.Name, slice); err != nil {
				return nil, err
			}
		default:
			// Third+ time — append to existing slice
			existing, err := result.GetMapKey(node.Name)
			if err != nil {
				return nil, err
			}
			if err := existing.Append(val); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// nodeToValue converts a single KDL node to a model.Value.
func nodeToValue(node *internal.Node) (*model.Value, error) {
	hasArgs := len(node.Arguments) > 0
	hasProps := len(node.Properties) > 0
	hasChildren := len(node.Children) > 0

	// Node with no args, no props, no children → null
	if !hasArgs && !hasProps && !hasChildren {
		return model.NewNullValue(), nil
	}

	// Node with single argument, no properties, no children → scalar
	if len(node.Arguments) == 1 && !hasProps && !hasChildren {
		return kdlValueToModelValue(node.Arguments[0])
	}

	// Otherwise, build a map
	result := model.NewMapValue()

	// Multiple args → $args key
	if hasArgs {
		if len(node.Arguments) == 1 && (hasProps || hasChildren) {
			// Single arg with props/children: still use $args
			argsSlice := model.NewSliceValue()
			v, err := kdlValueToModelValue(node.Arguments[0])
			if err != nil {
				return nil, err
			}
			if err := argsSlice.Append(v); err != nil {
				return nil, err
			}
			if err := result.SetMapKey("$args", argsSlice); err != nil {
				return nil, err
			}
		} else if len(node.Arguments) > 1 {
			argsSlice := model.NewSliceValue()
			for _, arg := range node.Arguments {
				v, err := kdlValueToModelValue(arg)
				if err != nil {
					return nil, err
				}
				if err := argsSlice.Append(v); err != nil {
					return nil, err
				}
			}
			if err := result.SetMapKey("$args", argsSlice); err != nil {
				return nil, err
			}
		}
	}

	// Properties become keys
	for _, prop := range node.Properties {
		v, err := kdlValueToModelValue(prop.Value)
		if err != nil {
			return nil, err
		}
		if err := result.SetMapKey(prop.Key, v); err != nil {
			return nil, err
		}
	}

	// Children are merged into the map
	if hasChildren {
		childMap, err := nodesToValue(node.Children)
		if err != nil {
			return nil, err
		}
		// Merge child keys into result
		if err := childMap.RangeMap(func(key string, val *model.Value) error {
			return result.SetMapKey(key, val)
		}); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func kdlValueToModelValue(v *internal.Value) (*model.Value, error) {
	switch val := v.Value.(type) {
	case string:
		return model.NewStringValue(val), nil
	case int64:
		return model.NewIntValue(val), nil
	case float64:
		return model.NewFloatValue(val), nil
	case bool:
		return model.NewBoolValue(val), nil
	case nil:
		return model.NewNullValue(), nil
	default:
		return nil, fmt.Errorf("kdl: unsupported value type %T", val)
	}
}
