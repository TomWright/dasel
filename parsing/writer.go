package parsing

import (
	"bytes"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

var writers = map[Format]NewWriterFn{}

type WriterOptions struct {
	Compact bool
	Indent  string
	Ext     map[string]string
}

// DefaultWriterOptions returns the default writer options.
func DefaultWriterOptions() WriterOptions {
	return WriterOptions{
		Compact: false,
		Indent:  "  ",
		Ext:     make(map[string]string),
	}
}

// Writer writes a value to a byte slice.
type Writer interface {
	// Write writes a value to a byte slice.
	Write(*model.Value) ([]byte, error)
}

// NewWriterFn is a function that creates a new writer.
type NewWriterFn func(options WriterOptions) (Writer, error)

// RegisterWriter registers a new writer for the format.
func RegisterWriter(format Format, fn NewWriterFn) {
	writers[format] = fn
}

// DocumentSeparator is an interface that can be implemented by writers to allow for custom document separators.
type DocumentSeparator interface {
	// Separator returns the document separator.
	Separator() []byte
}

// MultiDocumentWriter is a writer that can write multiple documents.
func MultiDocumentWriter(w Writer) Writer {
	return &multiDocumentWriter{w: w}
}

type multiDocumentWriter struct {
	w Writer
}

// Write writes a value to a byte slice.
func (w *multiDocumentWriter) Write(value *model.Value) ([]byte, error) {
	if value.IsBranch() || value.IsSpread() {
		buf := new(bytes.Buffer)

		documentSeparator := []byte("\n")
		if ds, ok := w.w.(DocumentSeparator); ok {
			documentSeparator = ds.Separator()
		}

		totalDocuments, err := value.SliceLen()
		if err != nil {
			return nil, fmt.Errorf("failed to get document length: %w", err)
		}

		if err := value.RangeSlice(func(i int, v *model.Value) error {
			docBytes, err := w.w.Write(v)
			if err != nil {
				return fmt.Errorf("failed to write document %d: %w", i, err)
			}
			buf.Write(docBytes)

			if i < totalDocuments-1 {
				buf.Write(documentSeparator)
			}

			return nil
		}); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
	return w.w.Write(value)
}
