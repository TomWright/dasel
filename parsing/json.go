package parsing

import (
	"encoding/json"

	"github.com/tomwright/dasel/v3/model"
)

func NewJSONReader() (Reader, error) {
	return &jsonReader{}, nil
}

func NewJSONWriter() (Writer, error) {
	return &jsonWriter{}, nil
}

type jsonReader struct{}

func (j *jsonReader) Read(data []byte) (*model.Value, error) {
	var unmarshalled any
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return nil, err
	}
	return model.NewValue(&unmarshalled), nil
}

type jsonWriter struct{}

func (j *jsonWriter) Write(value *model.Value) ([]byte, error) {
	return json.Marshal(value.Interface())
}
