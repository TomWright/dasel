package storage

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// YAMLParser is a Parser implementation to handle yaml files.
type YAMLParser struct {
}

// YAMLSingleDocument represents a decoded single-document YAML file.
type YAMLSingleDocument struct {
	originalRequired
	Value interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *YAMLSingleDocument) RealValue() interface{} {
	return d.Value
}

// YAMLMultiDocument represents a decoded multi-document YAML file.
type YAMLMultiDocument struct {
	originalRequired
	Values []interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *YAMLMultiDocument) RealValue() interface{} {
	return d.Values
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *YAMLParser) FromBytes(byteData []byte) (RealValue, error) {
	res := make([]interface{}, 0)

	decoder := yaml.NewDecoder(bytes.NewBuffer(byteData))

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
		return &YAMLSingleDocument{Value: res[0]}, nil
	default:
		return &YAMLMultiDocument{Values: res}, nil
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *YAMLParser) ToBytes(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()

	switch v := value.(type) {
	case *YAMLSingleDocument:
		if err := encoder.Encode(v.Value); err != nil {
			return nil, fmt.Errorf("could not encode single document: %w", err)
		}
	case *YAMLMultiDocument:
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
