package xml

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newXMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &xmlWriter{
		options: options,
	}, nil
}

type xmlWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *xmlWriter) Write(value *model.Value) ([]byte, error) {
	return nil, nil
}
