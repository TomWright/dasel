package dencoding

import (
	"bytes"
	"encoding/json"
	"io"
)

// lastOptions contains the options that the last JSONEncoder was created with.
// Find a better way of passing this information into nested MarshalJSON calls.
var lastOptions []JSONEncoderOption

// JSONEncoder wraps a standard json encoder to implement custom ordering logic.
type JSONEncoder struct {
	encoder *json.Encoder
}

// NewJSONEncoder returns a new dencoding JSONEncoder.
func NewJSONEncoder(w io.Writer, options ...JSONEncoderOption) *JSONEncoder {
	jsonEncoder := json.NewEncoder(w)
	encoder := &JSONEncoder{
		encoder: jsonEncoder,
	}
	for _, o := range options {
		o.ApplyEncoder(encoder)
	}
	lastOptions = options
	return encoder
}

// Encode encodes the given value and writes the encodes bytes to the stream.
func (encoder *JSONEncoder) Encode(v any) error {
	// We rely on Map.MarshalJSON to ensure ordering.
	return encoder.encoder.Encode(v)
}

// Close cleans up the encoder.
func (encoder *JSONEncoder) Close() error {
	return nil
}

// JSONEscapeHTML enables or disables html escaping when encoding JSON.
func JSONEscapeHTML(escape bool) JSONEncoderOption {
	return jsonEncodeHTMLOption{escapeHTML: escape}
}

type jsonEncodeHTMLOption struct {
	escapeHTML bool
}

func (option jsonEncodeHTMLOption) ApplyEncoder(encoder *JSONEncoder) {
	encoder.encoder.SetEscapeHTML(option.escapeHTML)
}

// JSONEncodeIndent sets the indentation when encoding JSON.
func JSONEncodeIndent(prefix string, indent string) JSONEncoderOption {
	return jsonEncodeIndent{prefix: prefix, indent: indent}
}

type jsonEncodeIndent struct {
	prefix string
	indent string
}

func (option jsonEncodeIndent) ApplyEncoder(encoder *JSONEncoder) {
	encoder.encoder.SetIndent(option.prefix, option.indent)
}

// MarshalJSON JSON encodes the map and returns the bytes.
// This maintains ordering.
func (m *Map) MarshalJSON() ([]byte, error) {

	buf := new(bytes.Buffer)
	buf.Write([]byte(`{`))
	encoder := NewJSONEncoder(buf, lastOptions...)
	for i, key := range m.keys {
		last := i == len(m.keys)-1

		if err := encoder.Encode(key); err != nil {
			return nil, err
		}
		buf.Write([]byte(`:`))
		if err := encoder.Encode(m.data[key]); err != nil {
			return nil, err
		}
		if !last {
			buf.Write([]byte(`,`))
		}
	}
	buf.Write([]byte(`}`))
	return buf.Bytes(), nil
}
