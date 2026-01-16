package yaml

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"go.yaml.in/yaml/v4"
)

var _ parsing.Reader = (*yamlReader)(nil)

func newYAMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &yamlReader{}, nil
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
			if v == nil {
				continue
			}
			if err := slice.Append(v.value); err != nil {
				return nil, err
			}
		}
		return slice, nil
	}
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
		yv.value.SetMetadataValue("yaml-alias", value.Value)
	}
	return nil
}
