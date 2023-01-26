package dencoding

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/pelletier/go-toml/v2/unstable"
	"io"
)

// TOMLDecoder wraps a standard toml encoder to implement custom ordering logic.
type TOMLDecoder struct {
	reader io.Reader
	p      *unstable.Parser
}

// NewTOMLDecoder returns a new dencoding TOMLDecoder.
func NewTOMLDecoder(r io.Reader, options ...TOMLDecoderOption) *TOMLDecoder {
	decoder := &TOMLDecoder{
		reader: r,
	}
	for _, o := range options {
		o.ApplyDecoder(decoder)
	}
	return decoder
}

// Decode decodes the next item found in the decoder and writes it to v.
func (decoder *TOMLDecoder) Decode(v any) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return io.EOF
	}
	return toml.Unmarshal(data, v)
}
