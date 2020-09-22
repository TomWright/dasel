package storage

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type unknownParserErr struct {
	parser string
}

func (e unknownParserErr) Error() string {
	return fmt.Sprintf("unknown parser: %s", e.parser)
}

type Parser interface {
	FromBytes(byteData []byte) (interface{}, error)
	LoadFromFile(filename string) (interface{}, error)
}

func FromString(parser string) (Parser, error) {
	switch parser {
	case "yaml":
		return &YAMLParser{}, nil
	default:
		return nil, &unknownParserErr{parser: parser}
	}
}

type YAMLParser struct {
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *YAMLParser) FromBytes(byteData []byte) (interface{}, error) {
	var data interface{}
	if err := yaml.Unmarshal(byteData, &data); err != nil {
		return data, fmt.Errorf("could not unmarshal config data: %w", err)
	}
	return data, nil
}

// LoadFromFile loads Data from the given file.
func (p *YAMLParser) LoadFromFile(filename string) (interface{}, error) {
	byteData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	return p.FromBytes(byteData)
}
