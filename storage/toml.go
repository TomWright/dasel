package storage

import (
	"bytes"
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/tomwright/dasel"
)

func init() {
	registerReadParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
	registerWriteParser([]string{"toml"}, []string{".toml"}, &TOMLParser{})
}

// TOMLParser is a Parser implementation to handle toml files.
type TOMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *TOMLParser) FromBytes(byteData []byte) (dasel.Value, error) {
	var data interface{}
	if err := toml.Unmarshal(byteData, &data); err != nil {
		return dasel.Value{}, fmt.Errorf("could not unmarshal data: %w", err)
	}
	return dasel.ValueOf(data).WithMetadata("isSingleDocument", true), nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *TOMLParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)

	enc := toml.NewEncoder(buf)

	colourise := false

	for _, o := range options {
		switch o.Key {
		case OptionIndent:
			if indent, ok := o.Value.(string); ok {
				enc.Indentation(indent)
			}
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				colourise = value
			}
		}
	}

	switch {
	case value.Metadata("isSingleDocument") == true:
		if err := enc.Encode(value.Interface()); err != nil {
			if err.Error() == "Only a struct or map can be marshaled to TOML" {
				buf.Write([]byte(fmt.Sprintf("%v\n", value.Interface())))
			} else {
				return nil, err
			}
		}
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			field := value.Index(i)
			if err := enc.Encode(field.Interface()); err != nil {
				if err.Error() == "Only a struct or map can be marshaled to TOML" {
					buf.Write([]byte(fmt.Sprintf("%v\n", field.Interface())))
				} else {
					return nil, err
				}
			}
		}
	default:
		if err := enc.Encode(value.Interface()); err != nil {
			if err.Error() == "Only a struct or map can be marshaled to TOML" {
				buf.Write([]byte(fmt.Sprintf("%v\n", value.Interface())))
			} else {
				return nil, err
			}
		}
	}

	if colourise {
		if err := ColouriseBuffer(buf, "toml"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buf.Bytes(), nil
}
