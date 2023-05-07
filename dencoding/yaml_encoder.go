package dencoding

import (
	"github.com/tomwright/dasel/v2/util"
	"gopkg.in/yaml.v3"
	"io"
)

// YAMLEncoder wraps a standard yaml encoder to implement custom ordering logic.
type YAMLEncoder struct {
	encoder *yaml.Encoder
}

// NewYAMLEncoder returns a new dencoding YAMLEncoder.
func NewYAMLEncoder(w io.Writer, options ...YAMLEncoderOption) *YAMLEncoder {
	yamlEncoder := yaml.NewEncoder(w)
	encoder := &YAMLEncoder{
		encoder: yamlEncoder,
	}
	for _, o := range options {
		o.ApplyEncoder(encoder)
	}
	return encoder
}

// Encode encodes the given value and writes the encodes bytes to the stream.
func (encoder *YAMLEncoder) Encode(v any) error {
	// We rely on Map.MarshalYAML to ensure ordering.
	return encoder.encoder.Encode(v)
}

// Close cleans up the encoder.
func (encoder *YAMLEncoder) Close() error {
	return encoder.encoder.Close()
}

// MarshalYAML YAML encodes the map and returns the bytes.
// This maintains ordering.
func (m *Map) MarshalYAML() (any, error) {
	return yamlOrderedMapToNode(m)
}

// YAMLEncodeIndent sets the indentation when encoding YAML.
func YAMLEncodeIndent(spaces int) YAMLEncoderOption {
	return yamlEncodeIndent{spaces: spaces}
}

type yamlEncodeIndent struct {
	spaces int
}

func (option yamlEncodeIndent) ApplyEncoder(encoder *YAMLEncoder) {
	encoder.encoder.SetIndent(option.spaces)
}

func yamlValueToNode(value any) (*yaml.Node, error) {
	switch v := value.(type) {
	case *Map:
		return yamlOrderedMapToNode(v)
	case []any:
		return yamlSliceToNode(v)
	default:
		return yamlScalarToNode(v)
	}
}

func yamlOrderedMapToNode(value *Map) (*yaml.Node, error) {
	mapNode := &yaml.Node{
		Kind:    yaml.MappingNode,
		Style:   yaml.TaggedStyle & yaml.DoubleQuotedStyle & yaml.SingleQuotedStyle & yaml.LiteralStyle & yaml.FoldedStyle & yaml.FlowStyle,
		Content: make([]*yaml.Node, 0),
	}

	for _, key := range value.keys {
		keyNode, err := yamlValueToNode(key)
		if err != nil {
			return nil, err
		}
		valueNode, err := yamlValueToNode(value.data[key])
		if err != nil {
			return nil, err
		}
		mapNode.Content = append(mapNode.Content, keyNode, valueNode)
	}

	return mapNode, nil
}

func yamlSliceToNode(value []any) (*yaml.Node, error) {
	node := &yaml.Node{
		Kind:    yaml.SequenceNode,
		Content: make([]*yaml.Node, len(value)),
	}

	for i, v := range value {
		indexNode, err := yamlValueToNode(v)
		if err != nil {
			return nil, err
		}
		node.Content[i] = indexNode
	}

	return node, nil
}

func yamlScalarToNode(value any) (*yaml.Node, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: util.ToString(value),
	}, nil
}
