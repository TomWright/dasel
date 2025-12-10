package yaml

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"gopkg.in/yaml.v3"
)

var _ parsing.Writer = (*yamlWriter)(nil)

func newYAMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &yamlWriter{}, nil
}

type yamlWriter struct{}

func (j *yamlWriter) Separator() []byte {
	return []byte("---\n")
}

// Write writes a value to a byte slice.
func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	yv := &yamlValue{value: value}
	res, err := yv.ToNode()
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(res)
}

func (yv *yamlValue) ToNode() (*yaml.Node, error) {
	res := &yaml.Node{}

	yamlAlias, ok := yv.value.Metadata["yaml-alias"].(string)
	if ok {
		//res.Kind = yaml.ScalarNode
		res.Kind = yaml.AliasNode
		res.Value = yamlAlias
		//res.Alias = &yaml.Node{
		//	Kind:  yaml.ScalarNode,
		//	Value: yamlAlias,
		//}
		return res, nil
	}

	switch yv.value.Type() {
	case model.TypeString:
		v, err := yv.value.StringValue()
		if err != nil {
			return nil, err
		}
		res.Kind = yaml.ScalarNode
		res.Value = v
		res.Tag = "!!str"
	case model.TypeBool:
		v, err := yv.value.BoolValue()
		if err != nil {
			return nil, err
		}
		res.Kind = yaml.ScalarNode
		res.Value = fmt.Sprintf("%t", v)
		res.Tag = "!!bool"
	case model.TypeInt:
		v, err := yv.value.IntValue()
		if err != nil {
			return nil, err
		}
		res.Kind = yaml.ScalarNode
		res.Value = fmt.Sprintf("%d", v)
		res.Tag = "!!int"
	case model.TypeFloat:
		v, err := yv.value.FloatValue()
		if err != nil {
			return nil, err
		}
		res.Kind = yaml.ScalarNode
		res.Value = fmt.Sprintf("%g", v)
		res.Tag = "!!float"
	case model.TypeMap:
		res.Kind = yaml.MappingNode
		if err := yv.value.RangeMap(func(key string, val *model.Value) error {
			keyNode := &yamlValue{value: model.NewStringValue(key)}
			valNode := &yamlValue{value: val}

			marshalledKey, err := keyNode.ToNode()
			if err != nil {
				return err
			}
			marshalledVal, err := valNode.ToNode()
			if err != nil {
				return err
			}

			res.Content = append(res.Content, marshalledKey)
			res.Content = append(res.Content, marshalledVal)

			return nil
		}); err != nil {
			return nil, err
		}
	case model.TypeSlice:
		res.Kind = yaml.SequenceNode
		if err := yv.value.RangeSlice(func(i int, val *model.Value) error {
			valNode := &yamlValue{value: val}
			marshalledVal, err := valNode.ToNode()
			if err != nil {
				return err
			}
			res.Content = append(res.Content, marshalledVal)
			return nil
		}); err != nil {
			return nil, err
		}
	case model.TypeNull:
		res.Kind = yaml.DocumentNode
	case model.TypeUnknown:
		return nil, fmt.Errorf("unknown type: %s", yv.value.Type())
	}

	return res, nil
}

func (yv *yamlValue) MarshalYAML() (any, error) {
	res, err := yv.ToNode()
	if err != nil {
		return nil, err
	}
	return res, nil
}
