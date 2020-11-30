package storage

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

func init() {
	registerReadParser([]string{"yaml", "yml"}, []string{".yaml", ".yml"}, &YAMLParser{})
	registerWriteParser([]string{"yaml", "yml"}, []string{".yaml", ".yml"}, &YAMLParser{})
}

// YAMLParser is a Parser implementation to handle yaml files.
type YAMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *YAMLParser) FromBytes(byteData []byte) (interface{}, error) {
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
		return &BasicSingleDocument{Value: res[0]}, nil
	default:
		return &BasicMultiDocument{Values: res}, nil
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *YAMLParser) ToBytes(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()

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
	return buffer.Bytes(), nil
}
