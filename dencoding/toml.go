package dencoding

// TOMLEncoderOption is identifies an option that can be applied to a TOML encoder.
type TOMLEncoderOption interface {
	ApplyEncoder(encoder *TOMLEncoder)
}

// TOMLDecoderOption is identifies an option that can be applied to a TOML decoder.
type TOMLDecoderOption interface {
	ApplyDecoder(decoder *TOMLDecoder)
}
