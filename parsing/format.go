package parsing

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
	TOML Format = "toml"
)

func (f Format) String() string {
	return string(f)
}

type Reader interface {
	Read([]byte) (*model.Value, error)
}

type Writer interface {
	Write(*model.Value) ([]byte, error)
}

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
