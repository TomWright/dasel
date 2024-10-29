package yaml

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"gopkg.in/yaml.v3"
)

// YAML represents the YAML file format.
const YAML parsing.Format = "yaml"

var _ parsing.Reader = (*yamlReader)(nil)
var _ parsing.Writer = (*yamlWriter)(nil)

func init() {
	parsing.RegisterReader(YAML, newYAMLReader)
	parsing.RegisterWriter(YAML, newYAMLWriter)
}

func newYAMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &yamlReader{}, nil
}

func newYAMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &yamlWriter{}, nil
}

type yamlReader struct{}

// Read reads a value from a byte slice.
func (j *yamlReader) Read(data []byte) (*model.Value, error) {
	d := yaml.NewDecoder(bytes.NewReader(data))
	res := make([]*yamlValue, 0)
	for {
		unmarshalled := &yamlValue{}
		if err := d.Decode(&unmarshalled); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		res = append(res, unmarshalled)
	}

	switch len(res) {
	case 0:
		return model.NewNullValue(), nil
	case 1:
		return res[0].value, nil
	default:
		slice := model.NewSliceValue()
		slice.MarkAsBranch()
		for _, v := range res {
			if err := slice.Append(v.value); err != nil {
				return nil, err
			}
		}
		return slice, nil
	}
}

type yamlWriter struct{}

// Write writes a value to a byte slice.
func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	if value.IsBranch() {
		res := make([]byte, 0)
		sliceLen, err := value.SliceLen()
		if err != nil {
			return nil, err
		}
		if err := value.RangeSlice(func(i int, val *model.Value) error {
			yv := &yamlValue{value: val}
			marshalled, err := yaml.Marshal(yv)
			if err != nil {
				return err
			}
			res = append(res, marshalled...)
			if i < sliceLen-1 {
				res = append(res, []byte("---\n")...)
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return res, nil
	}

	yv := &yamlValue{value: value}
	res, err := yv.ToNode()
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(res)
}

type yamlValue struct {
	node  *yaml.Node
	value *model.Value
}

func (yv *yamlValue) UnmarshalYAML(value *yaml.Node) error {
	yv.node = value
	switch value.Kind {
	case yaml.ScalarNode:
		switch value.Tag {
		case "!!bool":
			yv.value = model.NewBoolValue(value.Value == "true")
		case "!!int":
			i, err := strconv.Atoi(value.Value)
			if err != nil {
				return err
			}
			yv.value = model.NewIntValue(int64(i))
		case "!!float":
			f, err := strconv.ParseFloat(value.Value, 64)
			if err != nil {
				return err
			}
			yv.value = model.NewFloatValue(f)
		default:
			yv.value = model.NewStringValue(value.Value)
		}
	case yaml.DocumentNode:
		yv.value = model.NewNullValue()
	case yaml.SequenceNode:
		res := model.NewSliceValue()
		for _, item := range value.Content {
			newItem := &yamlValue{}
			if err := newItem.UnmarshalYAML(item); err != nil {
				return err
			}
			if err := res.Append(newItem.value); err != nil {
				return err
			}
		}
		yv.value = res
	case yaml.MappingNode:
		res := model.NewMapValue()
		for i := 0; i < len(value.Content); i += 2 {
			key := value.Content[i]
			val := value.Content[i+1]

			newKey := &yamlValue{}
			if err := newKey.UnmarshalYAML(key); err != nil {
				return err
			}

			newVal := &yamlValue{}
			if err := newVal.UnmarshalYAML(val); err != nil {
				return err
			}

			keyStr, err := newKey.value.StringValue()
			if err != nil {
				return fmt.Errorf("keys are expected to be strings: %w", err)
			}

			if err := res.SetMapKey(keyStr, newVal.value); err != nil {
				return err
			}
		}
		yv.value = res
	case yaml.AliasNode:
		newVal := &yamlValue{}
		if err := newVal.UnmarshalYAML(value.Alias); err != nil {
			return err
		}
		yv.value = newVal.value
		yv.value.Metadata["yaml-alias"] = value.Value
	}
	return nil
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
