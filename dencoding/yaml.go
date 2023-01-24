package dencoding

const (
	yamlTagString = "!!str"
	yamlTagMap    = "!!map"
	yamlTagArray  = "!!seq"
	yamlTagNull   = "!!null"
	yamlTagBinary = "!!binary"
	yamlTagBool   = "!!bool"
	yamlTagInt    = "!!int"
	yamlTagFloat  = "!!float"
)

// YAMLEncoderOption is identifies an option that can be applied to a YAML encoder.
type YAMLEncoderOption interface {
	ApplyEncoder(encoder *YAMLEncoder)
}

// YAMLDecoderOption is identifies an option that can be applied to a YAML decoder.
type YAMLDecoderOption interface {
	ApplyDecoder(decoder *YAMLDecoder)
}
