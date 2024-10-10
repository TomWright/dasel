package parsing

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// Format represents a file format.
type Format string

// Supported file formats.
const (
	JSON Format = "json"
	YAML Format = "yaml"
	TOML Format = "toml"
)

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// Reader reads a value from a byte slice.
type Reader interface {
	// Read reads a value from a byte slice.
	Read([]byte) (*model.Value, error)
}

// Writer writes a value to a byte slice.
type Writer interface {
	// Write writes a value to a byte slice.
	Write(*model.Value) ([]byte, error)
}

// NewReader creates a new reader for the specified format.
func NewReader(format Format) (Reader, error) {
	switch format {
	case JSON:
		return NewJSONReader()
	case YAML:
		return NewYAMLReader()
	case TOML:
		return NewTOMLReader()
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}

// NewWriter creates a new writer for the specified format.
func NewWriter(format Format) (Writer, error) {
	switch format {
	case JSON:
		return NewJSONWriter()
	case YAML:
		return NewYAMLWriter()
	case TOML:
		return NewTOMLWriter()
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}
