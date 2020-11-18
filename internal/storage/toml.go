package storage

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

// TOMLParser is a Parser implementation to handle toml files.
type TOMLParser struct {
	data interface{}
}

//RealValue implements interface RealValue for TOMLParser
func (p *TOMLParser) RealValue() interface{} {
	return p.data
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *TOMLParser) FromBytes(byteData []byte) (RealValue, error) {
	var data interface{}
	if err := toml.Unmarshal(byteData, &data); err != nil {
		return &TOMLParser{data: data}, fmt.Errorf("could not unmarshal config data: %w", err)
	}
	return &TOMLParser{data: data}, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value interface{}) ([]byte, error) {
	byteData, err := toml.Marshal(value)
	if err != nil && err.Error() == "Only a struct or map can be marshaled to TOML" {
		return []byte(fmt.Sprintf("%v\n", value)), nil
	}
	return byteData, err
}
