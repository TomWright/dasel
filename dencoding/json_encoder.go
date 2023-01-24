package dencoding

import (
	"encoding/json"
	"io"
)

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
	return encoder
}

// Encode encodes the given value and writes the encodes bytes to the stream.
func (encoder *JSONEncoder) Encode(v any) error {
	return encoder.encoder.Encode(v)
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
