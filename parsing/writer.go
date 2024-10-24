package parsing

import "github.com/tomwright/dasel/v3/model"

var writers = map[Format]NewWriterFn{}

type WriterOptions struct {
	Compact bool
	Indent  string
}

// DefaultWriterOptions returns the default writer options.
func DefaultWriterOptions() WriterOptions {
	return WriterOptions{
		Compact: false,
		Indent:  "  ",
	}
}

// Writer writes a value to a byte slice.
type Writer interface {
	// Write writes a value to a byte slice.
	Write(*model.Value) ([]byte, error)
}

// NewWriterFn is a function that creates a new writer.
type NewWriterFn func(options WriterOptions) (Writer, error)

// RegisterWriter registers a new writer for the format.
func RegisterWriter(format Format, fn NewWriterFn) {
	writers[format] = fn
}
