package storage

import (
	"fmt"
	"io"
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
	case ".toml":
		return &TOMLParser{}, nil
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
	case "toml":
		return &TOMLParser{}, nil
	default:
		return nil, &UnknownParserErr{Parser: parser}
	}
}

// LoadFromFile loads data from the given file.
func LoadFromFile(filename string, p Parser) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	return Load(p, f)
}

// Load loads data from the given io.Reader.
func Load(p Parser, reader io.Reader) (interface{}, error) {
	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %w", err)
	}
	return p.FromBytes(byteData)
}

// Write writes the value to the given io.Writer.
func Write(p Parser, value interface{}, writer io.Writer) error {
	byteData, err := p.ToBytes(value)
	if err != nil {
		return fmt.Errorf("could not get byte data for file: %w", err)
	}
	byteData = append(byteData, []byte("\n")...)
	if _, err := writer.Write(byteData); err != nil {
		return fmt.Errorf("could not write data: %w", err)
	}
	return nil
}
