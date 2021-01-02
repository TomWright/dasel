package storage

import (
	"bytes"
	"fmt"
)

func init() {
	registerWriteParser([]string{"-", "plain"}, []string{}, &PlainParser{})
}

// PlainParser is a Parser implementation to handle yaml files.
type PlainParser struct {
}

// ErrPlainParserNotImplemented is returned when you try to use the PlainParser.FromBytes func.
var ErrPlainParserNotImplemented = fmt.Errorf("PlainParser.FromBytes not implemented")

// FromBytes returns some data that is represented by the given bytes.
func (p *PlainParser) FromBytes(byteData []byte) (interface{}, error) {
	return nil, ErrPlainParserNotImplemented
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *PlainParser) ToBytes(value interface{}, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)
	switch val := value.(type) {
	case SingleDocument:
		buf.Write([]byte(fmt.Sprintf("%v\n", val.Document())))
	case MultiDocument:
		for _, doc := range val.Documents() {
			buf.Write([]byte(fmt.Sprintf("%v\n", doc)))
		}
	default:
		buf.Write([]byte(fmt.Sprintf("%v\n", val)))
	}
	return buf.Bytes(), nil
}
