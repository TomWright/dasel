package parsing

import "github.com/tomwright/dasel/v3/model"

func NewYAMLReader() (Reader, error) {
	return &yamlReader{}, nil
}

func NewYAMLWriter() (Writer, error) {
	return &yamlWriter{}, nil
}

type yamlReader struct{}

func (j *yamlReader) Read(data []byte) (*model.Value, error) {
	panic("not implemented")
}

type yamlWriter struct{}

func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	panic("not implemented")
}
