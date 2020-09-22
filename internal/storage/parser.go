package storage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// UnknownParserErr is returned when an invalid parser name is given.
type UnknownParserErr struct {
	parser string
}

// Error returns the error message.
func (e UnknownParserErr) Error() string {
	return fmt.Sprintf("unknown parser: %s", e.parser)
}

// Parser can be used to load and save files from/to disk.
type Parser interface {
	FromBytes(byteData []byte) (interface{}, error)
}

// NewParserFromFilename returns a Parser from the given filename.
func NewParserFromFilename(filename string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".yaml", ".yml":
		return &YAMLParser{}, nil
	case ".json":
		return &JSONParser{}, nil
	default:
		return nil, &UnknownParserErr{parser: ext}
	}
}

// NewParserFromString returns a Parser from the given parser name.
func NewParserFromString(parser string) (Parser, error) {
	switch parser {
	case "yaml":
		return &YAMLParser{}, nil
	case "json":
		return &JSONParser{}, nil
	default:
		return nil, &UnknownParserErr{parser: parser}
	}
}

// LoadFromFile loads data from the given file.
func LoadFromFile(filename string, p Parser) (interface{}, error) {
	byteData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	return p.FromBytes(byteData)
}
