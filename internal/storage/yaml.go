package storage

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// YAMLParser is a Parser implementation to handle yaml files.
type YAMLParser struct {
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *YAMLParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := yaml.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal config data: %w", err)
	}
	return data, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *YAMLParser) ToBytes(value interface{}) ([]byte, error) {
	return yaml.Marshal(value)
}
