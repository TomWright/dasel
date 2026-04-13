package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"go.yaml.in/yaml/v4"
)

var _ parsing.Reader = (*yamlReader)(nil)

func newYAMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &yamlReader{
		maxExpansionDepth:  maxExpansionDepth,
		maxExpansionBudget: maxExpansionBudget,
	}, nil
}

type yamlReader struct {
	maxExpansionDepth  int
	maxExpansionBudget int
}

// ErrYamlExpansionDepthExceeded is returned when the maximum expansion depth is exceeded.
var ErrYamlExpansionDepthExceeded = errors.New("yaml expansion depth exceeded")

// ErrYamlExpansionBudgetExceeded is returned when the maximum expansion budget is exceeded.
var ErrYamlExpansionBudgetExceeded = errors.New("yaml expansion budget exceeded")

const maxExpansionDepth = 32
const maxExpansionBudget = 1000

// Read reads a value from a byte slice.
func (j *yamlReader) Read(data []byte) (*model.Value, error) {
	d := yaml.NewDecoder(bytes.NewReader(data))
	res := make([]*yamlValue, 0)
	for {
		expansionBudget := j.maxExpansionBudget
		unmarshalled := &yamlValue{
			expansionDepth:    0,
			maxExpansionDepth: j.maxExpansionDepth,
			expansionBudget:   &expansionBudget,
		}
		if err := d.Decode(&unmarshalled); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if unmarshalled == nil {
			expansionBudget := j.maxExpansionBudget
			unmarshalled = &yamlValue{
				node:              nil,
				value:             model.NewNullValue(),
				expansionDepth:    0,
				maxExpansionDepth: j.maxExpansionDepth,
				expansionBudget:   &expansionBudget,
			}
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
	if yv.expansionDepth > yv.maxExpansionDepth {
		return ErrYamlExpansionDepthExceeded
	}
	switch value.Kind {
	case yaml.ScalarNode:
		switch value.Tag {
		case "!!bool":
			yv.value = model.NewBoolValue(value.Value == "true")
		case "!!int":
			i, err := parseYAMLInt(value.Value)
			if err != nil {
				return err
			}
			yv.value = model.NewIntValue(i)
		case "!!float":
			f, err := strconv.ParseFloat(value.Value, 64)
			if err != nil {
				return err
			}
			yv.value = model.NewFloatValue(f)
		case "!!null":
			yv.value = model.NewNullValue()
		case "!!str":
			yv.value = model.NewStringValue(value.Value)
		default:
			yv.value = model.NewStringValue(value.Value)
		}
	case yaml.DocumentNode:
		yv.value = model.NewNullValue()
	case yaml.SequenceNode:
		res := model.NewSliceValue()
		for _, item := range value.Content {
			newItem := &yamlValue{
				expansionDepth:    yv.expansionDepth,
				maxExpansionDepth: yv.maxExpansionDepth,
				expansionBudget:   yv.expansionBudget,
			}
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

			newKey := &yamlValue{
				expansionDepth:    yv.expansionDepth,
				maxExpansionDepth: yv.maxExpansionDepth,
				expansionBudget:   yv.expansionBudget,
			}
			if err := newKey.UnmarshalYAML(key); err != nil {
				return err
			}

			newVal := &yamlValue{
				expansionDepth:    yv.expansionDepth,
				maxExpansionDepth: yv.maxExpansionDepth,
				expansionBudget:   yv.expansionBudget,
			}
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
		if yv.expansionBudget != nil {
			*yv.expansionBudget = *yv.expansionBudget - 1
			if *yv.expansionBudget < 0 {
				return ErrYamlExpansionBudgetExceeded
			}
		}
		newVal := &yamlValue{
			expansionDepth:    yv.expansionDepth + 1,
			maxExpansionDepth: yv.maxExpansionDepth,
			expansionBudget:   yv.expansionBudget,
		}
		if err := newVal.UnmarshalYAML(value.Alias); err != nil {
			return err
		}
		yv.value = newVal.value
		yv.value.SetMetadataValue("yaml-alias", value.Value)
	}
	return nil
}

func parseYAMLInt(s string) (int64, error) {
	// Strip leading sign for prefix detection.
	clean := s
	if len(clean) > 0 && (clean[0] == '+' || clean[0] == '-') {
		clean = clean[1:]
	}

	switch {
	case strings.HasPrefix(clean, "0x") || strings.HasPrefix(clean, "0X"):
		return strconv.ParseInt(s, 0, 64)
	case strings.HasPrefix(clean, "0o") || strings.HasPrefix(clean, "0O"):
		return strconv.ParseInt(s, 0, 64)
	case strings.HasPrefix(clean, "0b") || strings.HasPrefix(clean, "0B"):
		return strconv.ParseInt(s, 0, 64)
	default:
		// YAML 1.2 allows underscores in decimal integers (e.g. 1_000).
		// strconv.ParseInt with base 10 does not support underscores,
		// so we strip them before parsing.
		return strconv.ParseInt(strings.ReplaceAll(s, "_", ""), 10, 64)
	}
}