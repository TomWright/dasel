package parsing

import (
	"encoding/json"

	"github.com/tomwright/dasel/v3/model"
)

// NewJSONReader creates a new JSON reader.
func NewJSONReader() (Reader, error) {
	return &jsonReader{}, nil
}

// NewJSONWriter creates a new JSON writer.
func NewJSONWriter() (Writer, error) {
	return &jsonWriter{}, nil
}

type jsonReader struct{}

// Read reads a value from a byte slice.
func (j *jsonReader) Read(data []byte) (*model.Value, error) {
	var unmarshalled any
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return nil, err
	}
	return model.NewValue(&unmarshalled), nil
}

type jsonWriter struct{}

// Write writes a value to a byte slice.
func (j *jsonWriter) Write(value *model.Value) ([]byte, error) {
	res, err := json.Marshal(value.Interface())
	if err != nil {
		return nil, err
	}
	return append(res, []byte("\n")...), nil
}
