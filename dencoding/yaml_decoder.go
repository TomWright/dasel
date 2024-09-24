package dencoding

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/tomwright/dasel/v2/util"
	"gopkg.in/yaml.v3"
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
	case yaml.AliasNode:
		return decoder.getNodeValue(node.Alias)
	default:
		return nil, fmt.Errorf("unhandled node kind: %v", node.Kind)
	}
}

func (decoder *YAMLDecoder) getMappingNodeValue(node *yaml.Node) (any, error) {
	res := NewMap()

	content := make([]*yaml.Node, 0)
	content = append(content, node.Content...)

	var keyNode *yaml.Node
	var valueNode *yaml.Node
	for {
		if len(content) == 0 {
			break
		}

		keyNode, valueNode, content = content[0], content[1], content[2:]

		if keyNode.ShortTag() == yamlTagMerge {
			content = append(valueNode.Alias.Content, content...)
			continue
		}

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
		return strconv.ParseInt(node.Value, 0, 64)
	case yamlTagString:
		return node.Value, nil
	case yamlTagTimestamp:
		value, ok := parseTimestamp(node.Value)
		if !ok {
			return value, fmt.Errorf("could not parse timestamp: %v", node.Value)
		}
		return value, nil
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

// This is a subset of the formats allowed by the regular expression
// defined at http://yaml.org/type/timestamp.html.
var allowedTimestampFormats = []string{
	"2006-1-2T15:4:5.999999999Z07:00", // RCF3339Nano with short date fields.
	"2006-1-2t15:4:5.999999999Z07:00", // RFC3339Nano with short date fields and lower-case "t".
	"2006-1-2 15:4:5.999999999",       // space separated with no time zone
	"2006-1-2",                        // date only
	// Notable exception: time.Tokenize cannot handle: "2001-12-14 21:59:43.10 -5"
	// from the set of examples.
}

// parseTimestamp parses s as a timestamp string and
// returns the timestamp and reports whether it succeeded.
// Timestamp formats are defined at http://yaml.org/type/timestamp.html
// Copied from yaml.v3.
func parseTimestamp(s string) (time.Time, bool) {
	// TODO write code to check all the formats supported by
	// http://yaml.org/type/timestamp.html instead of using time.Tokenize.

	// Quick check: all date formats start with YYYY-.
	i := 0
	for ; i < len(s); i++ {
		if c := s[i]; c < '0' || c > '9' {
			break
		}
	}
	if i != 4 || i == len(s) || s[i] != '-' {
		return time.Time{}, false
	}
	for _, format := range allowedTimestampFormats {
		if t, err := time.Parse(format, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
