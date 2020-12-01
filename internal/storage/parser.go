package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var readParsersByExtension = map[string]ReadParser{}
var writeParsersByExtension = map[string]WriteParser{}
var readParsersByName = map[string]ReadParser{}
var writeParsersByName = map[string]WriteParser{}

func registerReadParser(names []string, extensions []string, parser ReadParser) {
	for _, n := range names {
		readParsersByName[n] = parser
	}
	for _, e := range extensions {
		readParsersByExtension[e] = parser
	}
}

func registerWriteParser(names []string, extensions []string, parser WriteParser) {
	for _, n := range names {
		writeParsersByName[n] = parser
	}
	for _, e := range extensions {
		writeParsersByExtension[e] = parser
	}
}

// UnknownParserErr is returned when an invalid parser name is given.
type UnknownParserErr struct {
	Parser string
}

// Error returns the error message.
func (e UnknownParserErr) Error() string {
	return fmt.Sprintf("unknown parser: %s", e.Parser)
}

// ReadParser can be used to convert bytes to data.
type ReadParser interface {
	// FromBytes returns some data that is represented by the given bytes.
	FromBytes(byteData []byte) (interface{}, error)
}

// WriteParser can be used to convert data to bytes.
type WriteParser interface {
	// ToBytes returns a slice of bytes that represents the given value.
	ToBytes(value interface{}) ([]byte, error)
}

// Parser can be used to load and save files from/to disk.
type Parser interface {
	ReadParser
	WriteParser
}

// NewReadParserFromFilename returns a ReadParser from the given filename.
func NewReadParserFromFilename(filename string) (ReadParser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	p, ok := readParsersByExtension[ext]
	if !ok {
		return nil, &UnknownParserErr{Parser: ext}
	}
	return p, nil
}

// NewReadParserFromString returns a ReadParser from the given parser name.
func NewReadParserFromString(parser string) (ReadParser, error) {
	p, ok := readParsersByName[parser]
	if !ok {
		return nil, &UnknownParserErr{Parser: parser}
	}
	return p, nil
}

// NewWriteParserFromFilename returns a WriteParser from the given filename.
func NewWriteParserFromFilename(filename string) (WriteParser, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	p, ok := writeParsersByExtension[ext]
	if !ok {
		return nil, &UnknownParserErr{Parser: ext}
	}
	return p, nil
}

// NewWriteParserFromString returns a WriteParser from the given parser name.
func NewWriteParserFromString(parser string) (WriteParser, error) {
	p, ok := writeParsersByName[parser]
	if !ok {
		return nil, &UnknownParserErr{Parser: parser}
	}
	return p, nil
}

// LoadFromFile loads data from the given file.
func LoadFromFile(filename string, p ReadParser) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	return Load(p, f)
}

// Load loads data from the given io.Reader.
func Load(p ReadParser, reader io.Reader) (interface{}, error) {
	byteData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read data: %w", err)
	}
	return p.FromBytes(byteData)
}

// Write writes the value to the given io.Writer.
func Write(p WriteParser, value interface{}, originalValue interface{}, writer io.Writer) error {
	switch typed := originalValue.(type) {
	case OriginalRequired:
		if typed.OriginalRequired() {
			value = originalValue
		}
	}
	byteData, err := p.ToBytes(value)
	if err != nil {
		return fmt.Errorf("could not get byte data for file: %w", err)
	}
	if _, err := writer.Write(byteData); err != nil {
		return fmt.Errorf("could not write data: %w", err)
	}
	return nil
}

// OriginalRequired can be used in conjunction with RealValue to allow parsers to be more intelligent
// with the data they read/write.
type OriginalRequired interface {
	// OriginalRequired tells dasel if the parser requires the original value when converting to bytes.
	OriginalRequired() bool
}

// RealValue can be used in conjunction with OriginalRequired to allow parsers to be more intelligent
// with the data they read/write.
type RealValue interface {
	// RealValue returns the real value that dasel should use when processing data.
	RealValue() interface{}
}

type originalRequired struct {
}

// OriginalRequired tells dasel if the parser requires the original value when converting to bytes.
func (d originalRequired) OriginalRequired() bool {
	return true
}

type SingleDocument interface {
	Document() interface{}
}

type MultiDocument interface {
	Documents() []interface{}
}

// BasicSingleDocument represents a single document file.
type BasicSingleDocument struct {
	originalRequired
	Value interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *BasicSingleDocument) RealValue() interface{} {
	return d.Value
}

// Document returns the document that should be written to output.
func (d *BasicSingleDocument) Document() interface{} {
	return d.Value
}

// BasicMultiDocument represents a multi-document file.
type BasicMultiDocument struct {
	originalRequired
	Values []interface{}
}

// RealValue returns the real value that dasel should use when processing data.
func (d *BasicMultiDocument) RealValue() interface{} {
	return d.Values
}

// Documents returns the documents that should be written to output.
func (d *BasicMultiDocument) Documents() []interface{} {
	return d.Values
}
