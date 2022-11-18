package storage

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel"
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
func (p *YAMLParser) FromBytes(byteData []byte) (dasel.Value, error) {
	res := make([]interface{}, 0)

	decoder := yaml.NewDecoder(bytes.NewBuffer(byteData))

docLoop:
	for {
		var docData interface{}
		if err := decoder.Decode(&docData); err != nil {
			if err == io.EOF {
				break docLoop
			}
			return dasel.Value{}, fmt.Errorf("could not unmarshal data: %w", err)
		}

		formattedDocData := cleanupYamlMapValue(docData)

		res = append(res, formattedDocData)
	}
	switch len(res) {
	case 0:
		return dasel.Value{}, nil
	case 1:
		return dasel.ValueOf(res[0]).WithMetadata("isSingleDocument", true), nil
	default:
		return dasel.ValueOf(res).WithMetadata("isMultiDocument", true), nil
	}
}

func cleanupYamlInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = cleanupYamlMapValue(v)
	}
	return res
}

func cleanupYamlInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprint(k)] = cleanupYamlMapValue(v)
	}
	return res
}

func cleanupYamlMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanupYamlInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupYamlInterfaceMap(v)
	case string:
		return v
	default:
		return v
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *YAMLParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buffer)
	defer encoder.Close()

	colourise := false

	for _, o := range options {
		switch o.Key {
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				colourise = value
			}
		}
	}

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
		if err := ColouriseBuffer(buffer, "yaml"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buffer.Bytes(), nil
}
