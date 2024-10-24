package yaml

import (
	"bytes"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"gopkg.in/yaml.v3"
)

// YAML represents the YAML file format.
const YAML parsing.Format = "yaml"

var _ parsing.Reader = (*yamlReader)(nil)
var _ parsing.Writer = (*yamlWriter)(nil)

func init() {
	parsing.RegisterReader(YAML, newYAMLReader)
	parsing.RegisterWriter(YAML, newYAMLWriter)
}

func newYAMLReader() (parsing.Reader, error) {
	return &yamlReader{}, nil
}

func newYAMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &yamlWriter{}, nil
}

type yamlReader struct{}

// Read reads a value from a byte slice.
func (j *yamlReader) Read(data []byte) (*model.Value, error) {
	d := yaml.NewDecoder(bytes.NewReader(data))
	var unmarshalled any
	if err := d.Decode(&unmarshalled); err != nil {
		return nil, err
	}
	return model.NewValue(&unmarshalled), nil
}

type yamlWriter struct{}

// Write writes a value to a byte slice.
func (j *yamlWriter) Write(value *model.Value) ([]byte, error) {
	buf := new(bytes.Buffer)
	e := yaml.NewEncoder(buf)
	if err := e.Encode(value.Interface()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
