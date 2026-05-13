package kdl

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/kdl/internal"
)

var _ parsing.Writer = (*kdlWriter)(nil)

func newKDLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &kdlWriter{options: options}, nil
}

type kdlWriter struct {
	options parsing.WriterOptions
}

func (w *kdlWriter) Write(value *model.Value) ([]byte, error) {
	doc, err := modelValueToDocument(value)
	if err != nil {
		return nil, fmt.Errorf("kdl: %w", err)
	}

	opts := internal.DefaultGenerateOptions()
	opts.Compact = w.options.Compact
	if w.options.Indent != "" {
		opts.Indent = w.options.Indent
	}
	if v, ok := w.options.Ext["kdl-version"]; ok {
		switch v {
		case "1":
			opts.Version = internal.Version1
		case "2":
			opts.Version = internal.Version2
		default:
			return nil, fmt.Errorf("kdl: unsupported output version %q (use 1 or 2)", v)
		}
	}

	result, err := internal.GenerateString(doc, opts)
	if err != nil {
		return nil, fmt.Errorf("kdl: %w", err)
	}

	return []byte(result), nil
}

func modelValueToDocument(value *model.Value) (*internal.Document, error) {
	if value.Type() != model.TypeMap {
		// Non-map root: wrap as a single node
		node, err := modelValueToNode("root", value)
		if err != nil {
			return nil, err
		}
		return &internal.Document{Nodes: []*internal.Node{node}}, nil
	}

	nodes, err := mapToNodes(value)
	if err != nil {
		return nil, err
	}
	return &internal.Document{Nodes: nodes}, nil
}

func mapToNodes(value *model.Value) ([]*internal.Node, error) {
	kvs, err := value.MapKeyValues()
	if err != nil {
		return nil, err
	}

	var nodes []*internal.Node
	for _, kv := range kvs {
		if kv.Value.Type() == model.TypeSlice {
			// Slice → duplicate nodes with the same name
			sliceNodes, err := sliceToNodes(kv.Key, kv.Value)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, sliceNodes...)
		} else {
			node, err := modelValueToNode(kv.Key, kv.Value)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func sliceToNodes(name string, value *model.Value) ([]*internal.Node, error) {
	length, err := value.SliceLen()
	if err != nil {
		return nil, err
	}

	var nodes []*internal.Node
	for i := 0; i < length; i++ {
		elem, err := value.GetSliceIndex(i)
		if err != nil {
			return nil, err
		}
		node, err := modelValueToNode(name, elem)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func modelValueToNode(name string, value *model.Value) (*internal.Node, error) {
	node := &internal.Node{Name: name}

	switch value.Type() {
	case model.TypeString:
		s, err := value.StringValue()
		if err != nil {
			return nil, err
		}
		node.Arguments = append(node.Arguments, &internal.Value{Value: s})

	case model.TypeInt:
		n, err := value.IntValue()
		if err != nil {
			return nil, err
		}
		node.Arguments = append(node.Arguments, &internal.Value{Value: n})

	case model.TypeFloat:
		f, err := value.FloatValue()
		if err != nil {
			return nil, err
		}
		node.Arguments = append(node.Arguments, &internal.Value{Value: f})

	case model.TypeBool:
		b, err := value.BoolValue()
		if err != nil {
			return nil, err
		}
		node.Arguments = append(node.Arguments, &internal.Value{Value: b})

	case model.TypeNull:
		// Node with no args/props/children represents null

	case model.TypeMap:
		// Check for $args key
		kvs, err := value.MapKeyValues()
		if err != nil {
			return nil, err
		}

		var childNodes []*internal.Node
		for _, kv := range kvs {
			if kv.Key == "$args" {
				// $args → node arguments
				if kv.Value.Type() == model.TypeSlice {
					length, err := kv.Value.SliceLen()
					if err != nil {
						return nil, err
					}
					for i := 0; i < length; i++ {
						elem, err := kv.Value.GetSliceIndex(i)
						if err != nil {
							return nil, err
						}
						kdlVal, err := modelToKDLValue(elem)
						if err != nil {
							return nil, err
						}
						node.Arguments = append(node.Arguments, kdlVal)
					}
				}
				continue
			}

			// Scalar values → properties, complex values → children
			if isScalarType(kv.Value.Type()) {
				kdlVal, err := modelToKDLValue(kv.Value)
				if err != nil {
					return nil, err
				}
				node.Properties = append(node.Properties, &internal.Property{
					Key:   kv.Key,
					Value: kdlVal,
				})
			} else {
				// Map or slice value → child nodes
				if kv.Value.Type() == model.TypeSlice {
					sliceNodes, err := sliceToNodes(kv.Key, kv.Value)
					if err != nil {
						return nil, err
					}
					childNodes = append(childNodes, sliceNodes...)
				} else {
					childNode, err := modelValueToNode(kv.Key, kv.Value)
					if err != nil {
						return nil, err
					}
					childNodes = append(childNodes, childNode)
				}
			}
		}

		if len(childNodes) > 0 {
			node.Children = childNodes
		}

	case model.TypeSlice:
		// Slice as direct value → arguments
		length, err := value.SliceLen()
		if err != nil {
			return nil, err
		}
		for i := 0; i < length; i++ {
			elem, err := value.GetSliceIndex(i)
			if err != nil {
				return nil, err
			}
			kdlVal, err := modelToKDLValue(elem)
			if err != nil {
				return nil, err
			}
			node.Arguments = append(node.Arguments, kdlVal)
		}
	}

	return node, nil
}

func modelToKDLValue(value *model.Value) (*internal.Value, error) {
	switch value.Type() {
	case model.TypeString:
		s, err := value.StringValue()
		if err != nil {
			return nil, err
		}
		return &internal.Value{Value: s}, nil
	case model.TypeInt:
		n, err := value.IntValue()
		if err != nil {
			return nil, err
		}
		return &internal.Value{Value: n}, nil
	case model.TypeFloat:
		f, err := value.FloatValue()
		if err != nil {
			return nil, err
		}
		return &internal.Value{Value: f}, nil
	case model.TypeBool:
		b, err := value.BoolValue()
		if err != nil {
			return nil, err
		}
		return &internal.Value{Value: b}, nil
	case model.TypeNull:
		return &internal.Value{Value: nil}, nil
	default:
		return nil, fmt.Errorf("cannot convert %s to KDL value", value.Type())
	}
}

func isScalarType(t model.Type) bool {
	switch t {
	case model.TypeString, model.TypeInt, model.TypeFloat, model.TypeBool, model.TypeNull:
		return true
	}
	return false
}
