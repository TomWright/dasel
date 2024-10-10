package parsing

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
)

// NewTOMLReader creates a new TOML reader.
func NewTOMLReader() (Reader, error) {
	return &tomlReader{}, nil
}

// NewTOMLWriter creates a new TOML writer.
func NewTOMLWriter() (Writer, error) {
	return &tomlWriter{}, nil
}

type tomlReader struct{}

// Read reads a value from a byte slice.
func (j *tomlReader) Read(data []byte) (*model.Value, error) {
	var unmarshalled any
	if err := toml.Unmarshal(data, &unmarshalled); err != nil {
		return nil, err
	}
	return model.NewValue(&unmarshalled), nil
}

type tomlWriter struct{}

// Write writes a value to a byte slice.
func (j *tomlWriter) Write(value *model.Value) ([]byte, error) {
	res, err := toml.Marshal(value.Interface())
	if err != nil {
		return nil, err
	}
	return append(res, []byte("\n")...), nil
}
