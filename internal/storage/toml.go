package storage

import (
	"bytes"
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
	return &BasicSingleDocument{
		Value: data,
	}, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := toml.NewEncoder(buf)

	switch d := value.(type) {
	case SingleDocument:
		if err := enc.Encode(d.Document()); err != nil {
			if err.Error() == "Only a struct or map can be marshaled to TOML" {
				buf.Write([]byte(fmt.Sprintf("%v\n", d)))
			} else {
				return nil, err
			}
		}
	case MultiDocument:
		for _, dd := range d.Documents() {
			if err := enc.Encode(dd); err != nil {
				if err.Error() == "Only a struct or map can be marshaled to TOML" {
					buf.Write([]byte(fmt.Sprintf("%v\n", dd)))
				} else {
					return nil, err
				}
			}
		}
	default:
		if err := enc.Encode(d); err != nil {
			if err.Error() == "Only a struct or map can be marshaled to TOML" {
				buf.Write([]byte(fmt.Sprintf("%v\n", d)))
			} else {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}
