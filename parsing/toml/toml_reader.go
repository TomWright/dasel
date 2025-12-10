package toml

import (
	"github.com/pelletier/go-toml/v2"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Reader = (*tomlReader)(nil)

func newTOMLReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &tomlReader{}, nil
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
