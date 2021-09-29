package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func init() {
	registerReadParser([]string{"json"}, []string{".json"}, &JSONParser{})
	registerWriteParser([]string{"json"}, []string{".json"}, &JSONParser{})
}

// JSONParser is a Parser implementation to handle yaml files.
type JSONParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *JSONParser) FromBytes(byteData []byte) (interface{}, error) {
	res := make([]interface{}, 0)

	decoder := json.NewDecoder(bytes.NewBuffer(byteData))

docLoop:
	for {
		var docData interface{}
		if err := decoder.Decode(&docData); err != nil {
			if err == io.EOF {
				break docLoop
			}
			return nil, fmt.Errorf("could not unmarshal data: %w", err)
		}
		res = append(res, docData)
	}
	switch len(res) {
	case 0:
		return nil, nil
	case 1:
		return &BasicSingleDocument{Value: res[0]}, nil
	default:
		return &BasicMultiDocument{Values: res}, nil
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *JSONParser) ToBytes(value interface{}, options ...ReadWriteOption) ([]byte, error) {
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

	switch v := value.(type) {
	case SingleDocument:
		if err := encoder.Encode(v.Document()); err != nil {
			return nil, fmt.Errorf("could not encode single document: %w", err)
		}
	case MultiDocument:
		for index, d := range v.Documents() {
			if err := encoder.Encode(d); err != nil {
				return nil, fmt.Errorf("could not encode multi document [%d]: %w", index, err)
			}
		}
	default:
		if err := encoder.Encode(v); err != nil {
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
