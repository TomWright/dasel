package dencoding

import (
	"fmt"
	"github.com/tomwright/dasel/v2/util"
	"gopkg.in/yaml.v3"
	"io"
	"reflect"
	"strconv"
)

// YAMLDecoder wraps a standard yaml encoder to implement custom ordering logic.
type YAMLDecoder struct {
	decoder *yaml.Decoder
}

// NewYAMLDecoder returns a new dencoding YAMLDecoder.
func NewYAMLDecoder(r io.Reader, options ...YAMLDecoderOption) *YAMLDecoder {
	yamlDecoder := yaml.NewDecoder(r)
	decoder := &YAMLDecoder{
		decoder: yamlDecoder,
	}
	for _, o := range options {
		o.ApplyDecoder(decoder)
	}
	return decoder
}

// Decode decodes the next item found in the decoder and writes it to v.
func (decoder *YAMLDecoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("invalid decode target: %s", reflect.TypeOf(v))
	}

	rve := rv.Elem()

	node, err := decoder.nextNode()
	if err != nil {
		return err
	}

	if node.Kind == yaml.DocumentNode && len(node.Content) == 1 && node.Content[0].ShortTag() == yamlTagNull {
		return io.EOF
	}

	val, err := decoder.getNodeValue(node)
	if err != nil {
		return err
	}

	rve.Set(reflect.ValueOf(val))
	return nil
}

func (decoder *YAMLDecoder) getNodeValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		return decoder.getNodeValue(node.Content[0])
	case yaml.MappingNode:
		return decoder.getMappingNodeValue(node)
	case yaml.SequenceNode:
		return decoder.getSequenceNodeValue(node)
	case yaml.ScalarNode:
		return decoder.getScalarNodeValue(node)
	default:
		return nil, fmt.Errorf("unhandled node kind: %v", node.Kind)
	}
}

func (decoder *YAMLDecoder) getMappingNodeValue(node *yaml.Node) (any, error) {
	res := NewMap()

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		keyValue, err := decoder.getNodeValue(keyNode)
		if err != nil {
			return nil, err
		}
		value, err := decoder.getNodeValue(valueNode)
		if err != nil {
			return nil, err
		}

		key := util.ToString(keyValue)

		res.Set(key, value)
	}

	return res, nil
}

func (decoder *YAMLDecoder) getSequenceNodeValue(node *yaml.Node) (any, error) {
	res := make([]any, len(node.Content))
	for k, n := range node.Content {
		val, err := decoder.getNodeValue(n)
		if err != nil {
			return nil, err
		}
		res[k] = val
	}
	return res, nil
}

func (decoder *YAMLDecoder) getScalarNodeValue(node *yaml.Node) (any, error) {
	switch node.ShortTag() {
	case yamlTagNull:
		return nil, nil
	case yamlTagBool:
		return node.Value == "true", nil
	case yamlTagFloat:
		return strconv.ParseFloat(node.Value, 64)
	case yamlTagInt:
		return strconv.ParseInt(node.Value, 10, 64)
	case yamlTagString:
		return node.Value, nil
	default:
		return nil, fmt.Errorf("unhandled scalar node tag: %v", node.ShortTag())
	}
}

func (decoder *YAMLDecoder) nextNode() (*yaml.Node, error) {
	var node yaml.Node
	if err := decoder.decoder.Decode(&node); err != nil {
		return nil, err
	}
	return &node, nil
}
