package parsing

import "github.com/tomwright/dasel/v3/model"

func NewTOMLReader() (Reader, error) {
	return &tomlReader{}, nil
}

func NewTOMLWriter() (Writer, error) {
	return &tomlWriter{}, nil
}

type tomlReader struct{}

func (j *tomlReader) Read(data []byte) (*model.Value, error) {
	panic("not implemented")
}

type tomlWriter struct{}

func (j *tomlWriter) Write(value *model.Value) ([]byte, error) {
	panic("not implemented")
}
