package storage

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

func init() {
	registerReadParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
	registerWriteParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
}

// TOMLParser is a Parser implementation to handle toml files.
type TOMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *TOMLParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := toml.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal data: %w", err)
	}
	return data, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value interface{}) ([]byte, error) {
	byteData, err := toml.Marshal(value)
	if err != nil && err.Error() == "Only a struct or map can be marshaled to TOML" {
		return []byte(fmt.Sprintf("%v\n", value)), nil
	}
	return byteData, err
}
