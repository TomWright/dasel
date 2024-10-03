package parsing

import "github.com/tomwright/dasel/v3/model"

// NewYAMLReader creates a new YAML reader.
func NewYAMLReader() (Reader, error) {
	return &yamlReader{}, nil
}

// NewYAMLWriter creates a new YAML writer.
func NewYAMLWriter() (Writer, error) {
	return &yamlWriter{}, nil
}

type yamlReader struct{}

// Read reads a value from a byte slice.
func (j *yamlReader) Read(data []byte) (*model.Value, error) {
	panic("not implemented")
}

type yamlWriter struct{}

// Write writes a value to a byte slice.
func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	panic("not implemented")
}
