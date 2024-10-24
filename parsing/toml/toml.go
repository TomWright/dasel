package toml

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// TODO : Implement using https://github.com/pelletier/go-toml/blob/v2/unstable/ast.go

// TOML represents the TOML file format.
const TOML parsing.Format = "toml"

var _ parsing.Reader = (*tomlReader)(nil)
var _ parsing.Writer = (*tomlWriter)(nil)

func init() {
	parsing.RegisterReader(TOML, newTOMLReader)
	parsing.RegisterWriter(TOML, newTOMLWriter)
}

func newTOMLReader() (parsing.Reader, error) {
	return &tomlReader{}, nil
}

func newTOMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
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
