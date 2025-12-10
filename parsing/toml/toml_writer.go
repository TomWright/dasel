package toml

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Writer = (*tomlWriter)(nil)

func newTOMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &tomlWriter{}, nil
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
