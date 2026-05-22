package yaml

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"unicode/utf8"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"go.yaml.in/yaml/v4"
)

var _ parsing.Writer = (*yamlWriter)(nil)

func newYAMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &yamlWriter{options: options}, nil
}

type yamlWriter struct {
	options parsing.WriterOptions
}

func (j *yamlWriter) Separator() []byte {
	return []byte("---\n")
}

// Write writes a value to a byte slice.
func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	yv := &yamlValue{value: value, compact: j.options.Compact}
	res, err := yv.ToNode()
	if err != nil {
		return nil, err
	}
	out, err := yaml.Marshal(res)
	if err != nil {
		return nil, err
	}
	return unescapeYAMLUnicode(out), nil
}

// unescapeYAMLUnicode replaces \UNNNNNNNN escape sequences (8 hex digits) with
// the corresponding UTF-8 bytes. The yaml library incorrectly escapes
// supplementary-plane characters (U+10000–U+10FFFF) because its isPrintable
// function does not handle 4-byte UTF-8 sequences.
func unescapeYAMLUnicode(data []byte) []byte {
	marker := []byte(`\U`)
	var result []byte
	for {
		idx := bytes.Index(data, marker)
		if idx == -1 {
			break
		}
		// Need exactly 8 hex digits after \U
		if idx+10 > len(data) {
			result = append(result, data[:idx+2]...)
			data = data[idx+2:]
			continue
		}
		hexBytes := data[idx+2 : idx+10]
		decoded, err := hex.DecodeString(string(hexBytes))
		if err != nil || len(decoded) != 4 {
			result = append(result, data[:idx+2]...)
			data = data[idx+2:]
			continue
		}
		r := rune(decoded[0])<<24 | rune(decoded[1])<<16 | rune(decoded[2])<<8 | rune(decoded[3])
		if r < 0x10000 || r > 0x10FFFF || !utf8.ValidRune(r) {
			result = append(result, data[:idx+10]...)
			data = data[idx+10:]
			continue
		}
		result = append(result, data[:idx]...)
		var buf [4]byte
		n := utf8.EncodeRune(buf[:], r)
		result = append(result, buf[:n]...)
		data = data[idx+10:]
	}
	result = append(result, data...)
	return result
}

func (yv *yamlValue) ToNode() (*yaml.Node, error) {
	res := &yaml.Node{}

	// TODO : Handle yaml aliases.
	//yamlAlias, ok := yv.value.Metadata["yaml-alias"].(string)
	//if ok {
	//res.Kind = yaml.AliasNode
	//res.Value = yamlAlias
	//return res, nil
	//}

	switch yv.value.Type() {
	case model.TypeString:
		v, err := yv.value.StringValue()
		if err != nil {
			return nil, err
		}
		res.Kind = yaml.ScalarNode
		res.Value = v
		res.Tag = "!!str"
		if styleVal, ok := yv.value.MetadataValue("yaml-style"); ok {
			if style, ok := styleVal.(yaml.Style); ok {
				res.Style = style
			}
		}
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
		if yv.compact {
			res.Style = yaml.FlowStyle
		}
		if err := yv.value.RangeMap(func(key string, val *model.Value) error {
			keyNode := &yamlValue{value: model.NewStringValue(key), compact: yv.compact}
			valNode := &yamlValue{value: val, compact: yv.compact}

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
		if yv.compact {
			res.Style = yaml.FlowStyle
		}
		if err := yv.value.RangeSlice(func(i int, val *model.Value) error {
			valNode := &yamlValue{value: val, compact: yv.compact}
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
		res.Kind = yaml.ScalarNode
		res.Value = "null"
		res.Tag = "!!null"
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
