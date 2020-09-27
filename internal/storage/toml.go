package storage

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

// TOMLParser is a Parser implementation to handle yaml files.
type TOMLParser struct {
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *TOMLParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := toml.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal config data: %w", err)
	}
	return data, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value interface{}) ([]byte, error) {
	return toml.Marshal(value)
}
