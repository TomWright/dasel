package storage

import (
	"fmt"
)

func init() {
	registerWriteParser([]string{"-", "plain"}, []string{}, &PlainParser{})
}

// PlainParser is a Parser implementation to handle yaml files.
type PlainParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *PlainParser) FromBytes(byteData []byte) (interface{}, error) {
	return nil, fmt.Errorf("PlainParser.FromBytes not implemented")
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *PlainParser) ToBytes(value interface{}) ([]byte, error) {
	return []byte(fmt.Sprintf("%v\n", value)), nil
}
