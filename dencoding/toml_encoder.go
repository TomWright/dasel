package dencoding

import (
	"github.com/pelletier/go-toml/v2"
	"io"
)

// TOMLEncoder wraps a standard toml encoder to implement custom ordering logic.
type TOMLEncoder struct {
	encoder *toml.Encoder
	writer  io.Writer
}

// NewTOMLEncoder returns a new dencoding TOMLEncoder.
func NewTOMLEncoder(w io.Writer, options ...TOMLEncoderOption) *TOMLEncoder {
	tomlEncoder := toml.NewEncoder(w)
	encoder := &TOMLEncoder{
		writer:  w,
		encoder: tomlEncoder,
	}
	for _, o := range options {
		o.ApplyEncoder(encoder)
	}
	return encoder
}

// Encode encodes the given value and writes the encodes bytes to the stream.
func (encoder *TOMLEncoder) Encode(v any) error {
	// No ordering is done here.
	adjusted := removeDencodingMap(v)
	return encoder.encoder.Encode(adjusted)
}

// Close cleans up the encoder.
func (encoder *TOMLEncoder) Close() error {
	return nil
}

func removeDencodingMap(value any) any {
	switch v := value.(type) {
	case []any:
		return removeDencodingMapFromArray(v)
	case map[string]any:
		return removeDencodingMapFromMap(v)
	case *Map:
		return removeDencodingMap(v.data)
	default:
		return v
	}
}

func removeDencodingMapFromArray(value []any) []any {
	for k, v := range value {
		value[k] = removeDencodingMap(v)
	}
	return value
}

func removeDencodingMapFromMap(value map[string]any) map[string]any {
	for k, v := range value {
		value[k] = removeDencodingMap(v)
	}
	return value
}

// TOMLIndentSymbol sets the indentation when encoding TOML.
func TOMLIndentSymbol(symbol string) TOMLEncoderOption {
	return tomlEncodeSymbol{symbol: symbol}
}

type tomlEncodeSymbol struct {
	symbol string
}

func (option tomlEncodeSymbol) ApplyEncoder(encoder *TOMLEncoder) {
	encoder.encoder.SetIndentSymbol(option.symbol)
}
