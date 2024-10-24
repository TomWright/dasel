package parsing

import (
	"fmt"
)

// Format represents a file format.
type Format string

// NewReader creates a new reader for the format.
func (f Format) NewReader() (Reader, error) {
	fn, ok := readers[f]
	if !ok {
		return nil, fmt.Errorf("unsupported reader file format: %s", f)
	}
	return fn()
}

// NewWriter creates a new writer for the format.
func (f Format) NewWriter(options WriterOptions) (Writer, error) {
	fn, ok := writers[f]
	if !ok {
		return nil, fmt.Errorf("unsupported writer file format: %s", f)
	}
	return fn(options)
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// RegisteredReaders returns a list of registered readers.
func RegisteredReaders() []Format {
	var formats []Format
	for format := range readers {
		formats = append(formats, format)
	}
	return formats
}

// RegisteredWriters returns a list of registered writers.
func RegisteredWriters() []Format {
	var formats []Format
	for format := range writers {
		formats = append(formats, format)
	}
	return formats
}
