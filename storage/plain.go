package storage

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/v2"
)

func init() {
	registerWriteParser([]string{"-", "plain"}, []string{}, &PlainParser{})
}

// PlainParser is a Parser implementation to handle plain files.
type PlainParser struct {
}

// ErrPlainParserNotImplemented is returned when you try to use the PlainParser.FromBytes func.
var ErrPlainParserNotImplemented = fmt.Errorf("PlainParser.FromBytes not implemented")

// FromBytes returns some data that is represented by the given bytes.
func (p *PlainParser) FromBytes(byteData []byte) (dasel.Value, error) {
	return dasel.Value{}, ErrPlainParserNotImplemented
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *PlainParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)

	switch {
	case value.Metadata("isSingleDocument") == true:
		buf.Write([]byte(fmt.Sprintf("%v\n", value.Interface())))
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			buf.Write([]byte(fmt.Sprintf("%v\n", value.Index(i).Interface())))
		}
	default:
		buf.Write([]byte(fmt.Sprintf("%v\n", value.Interface())))
	}

	return buf.Bytes(), nil
}
