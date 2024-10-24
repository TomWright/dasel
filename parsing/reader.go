package parsing

import "github.com/tomwright/dasel/v3/model"

var readers = map[Format]NewReaderFn{}

// Reader reads a value from a byte slice.
type Reader interface {
	// Read reads a value from a byte slice.
	Read([]byte) (*model.Value, error)
}

// NewReaderFn is a function that creates a new reader.
type NewReaderFn func() (Reader, error)

// RegisterReader registers a new reader for the format.
func RegisterReader(format Format, fn NewReaderFn) {
	readers[format] = fn
}
