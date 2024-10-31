package hcl

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newHCLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &hclWriter{}, nil
}

type hclWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *hclWriter) Write(value *model.Value) ([]byte, error) {
	return nil, nil
}
