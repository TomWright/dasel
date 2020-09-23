package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// UnknownParserErr is returned when an invalid parser name is given.
type UnknownParserErr struct {
	Parser string
}

// Error returns the error message.
func (e UnknownParserErr) Error() string {
	return fmt.Sprintf("unknown parser: %s", e.Parser)
}

// Parser can be used to load and save files from/to disk.
type Parser interface {
	// FromBytes returns some Data that is represented by the given bytes.
	FromBytes(byteData []byte) (interface{}, error)
	// ToBytes returns a slice of bytes that represents the given value.
	ToBytes(value interface{}) ([]byte, error)
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
		return nil, &UnknownParserErr{Parser: ext}
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
		return nil, &UnknownParserErr{Parser: parser}
	}
}

// LoadFromFile loads data from the given file.
func LoadFromFile(filename string, p Parser) (interface{}, error) {
	byteData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}
	return p.FromBytes(byteData)
}

// LoadFromStdin loads data from the given file.
func LoadFromStdin(p Parser) (interface{}, error) {
	byteData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("could not read stdin: %w", err)
	}
	return p.FromBytes(byteData)
}

// WriteToFile saves data to the given file.
func WriteToFile(filename string, p Parser, value interface{}) error {
	byteData, err := p.ToBytes(value)
	if err != nil {
		return fmt.Errorf("could not get byte data for file: %w", err)
	}
	if err := ioutil.WriteFile(filename, byteData, 0644); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}
	return nil
}

// WriteToStdout saves data to the given file.
func WriteToStdout(p Parser, value interface{}) error {
	byteData, err := p.ToBytes(value)
	if err != nil {
		return fmt.Errorf("could not get byte data for file: %w", err)
	}
	if _, err := os.Stdout.Write(byteData); err != nil {
		return fmt.Errorf("could not write to stdout: %w", err)
	}
	return nil
}
