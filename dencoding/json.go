package dencoding

import "encoding/json"

const (
	jsonOpenObject  = json.Delim('{')
	jsonCloseObject = json.Delim('}')
	jsonOpenArray   = json.Delim('[')
	jsonCloseArray  = json.Delim(']')
)

// JSONEncoderOption is identifies an option that can be applied to a JSON encoder.
type JSONEncoderOption interface {
	ApplyEncoder(encoder *JSONEncoder)
}

// JSONDecoderOption is identifies an option that can be applied to a JSON decoder.
type JSONDecoderOption interface {
	ApplyDecoder(decoder *JSONDecoder)
}
