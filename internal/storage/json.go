package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// JSONParser is a Parser implementation to handle yaml files.
type JSONParser struct {
}

// JSONSingleDocument represents a decoded single-document YAML file.
type JSONSingleDocument struct {
	originalRequired
	Value interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *JSONSingleDocument) RealValue() interface{} {
	return d.Value
}

// YAMLMultiDocument represents a decoded multi-document YAML file.
type JSONMultiDocument struct {
	originalRequired
	Values []interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *JSONMultiDocument) RealValue() interface{} {
	return d.Values
}

// FromBytes returns some Data that is represented by the given bytes.

func (p *JSONParser) FromBytes(byteData []byte) (RealValue, error) {

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
		return &JSONSingleDocument{Value: res[0]}, nil
	default:
		return &JSONMultiDocument{Values: res}, nil
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *JSONParser) ToBytes(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")

	switch v := value.(type) {
	case *JSONSingleDocument:
		if err := encoder.Encode(v.Value); err != nil {
			return nil, fmt.Errorf("could not encode single document: %w", err)
		}
	case *JSONMultiDocument:
		for index, d := range v.Values {
			if err := encoder.Encode(d); err != nil {
				return nil, fmt.Errorf("could not encode multi document [%d]: %w", index, err)
			}
		}
	default:
		if err := encoder.Encode(v); err != nil {
			return nil, fmt.Errorf("could not encode default document type: %w", err)
		}
	}
	return buffer.Bytes(), nil
}
