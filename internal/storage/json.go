package storage

import (
	"encoding/json"
	"fmt"
)

// JSONParser is a Parser implementation to handle yaml files.
type JSONParser struct {
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *JSONParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal config data: %w", err)
	}
	return data, nil
}
