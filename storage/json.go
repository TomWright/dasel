package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tomwright/dasel"
	"io"
)

func init() {
	registerReadParser([]string{"json"}, []string{".json"}, &JSONParser{})
	registerWriteParser([]string{"json"}, []string{".json"}, &JSONParser{})
}

// JSONParser is a Parser implementation to handle json files.
type JSONParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *JSONParser) FromBytes(byteData []byte) (dasel.Value, error) {
	res := make([]interface{}, 0)

	decoder := json.NewDecoder(bytes.NewBuffer(byteData))

docLoop:
	for {
		var docData interface{}
		if err := decoder.Decode(&docData); err != nil {
			if err == io.EOF {
				break docLoop
			}
			return dasel.Value{}, fmt.Errorf("could not unmarshal data: %w", err)
		}
		res = append(res, docData)
	}

	switch len(res) {
	case 0:
		return dasel.Value{}, nil
	case 1:
		return dasel.ValueOf(res[0]).
			WithMetadata("isSingleDocument", true), nil
	default:
		return dasel.ValueOf(res).
			WithMetadata("isMultiDocument", true), nil
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *JSONParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)

	indent := "  "
	prettyPrint := true
	colourise := false

	for _, o := range options {
		switch o.Key {
		case OptionIndent:
			if value, ok := o.Value.(string); ok {
				indent = value
			}
		case OptionPrettyPrint:
			if value, ok := o.Value.(bool); ok {
				prettyPrint = value
			}
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				colourise = value
			}
		case OptionEscapeHTML:
			if value, ok := o.Value.(bool); ok {
				encoder.SetEscapeHTML(value)
			}
		}
	}

	if !prettyPrint {
		indent = ""
	}
	encoder.SetIndent("", indent)

	switch {
	case value.Metadata("isSingleDocument") == true:
		if err := encoder.Encode(value.Interface()); err != nil {
			return nil, fmt.Errorf("could not encode single document: %w", err)
		}
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			if err := encoder.Encode(value.Index(i).Interface()); err != nil {
				return nil, fmt.Errorf("could not encode multi document [%d]: %w", i, err)
			}
		}
	default:
		if err := encoder.Encode(value.Interface()); err != nil {
			return nil, fmt.Errorf("could not encode default document type: %w", err)
		}
	}

	if colourise {
		if err := ColouriseBuffer(buffer, "json"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buffer.Bytes(), nil
}
