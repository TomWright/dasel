package storage

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/v2"
	"strings"

	"github.com/clbanning/mxj/v2"
	"golang.org/x/net/html/charset"
)

func init() {
	// Required for https://github.com/TomWright/dasel/issues/61
	mxj.XMLEscapeCharsDecoder(true)

	// Required for https://github.com/TomWright/dasel/issues/164
	mxj.XmlCharsetReader = charset.NewReaderLabel

	registerReadParser([]string{"xml"}, []string{".xml"}, &XMLParser{})
	registerWriteParser([]string{"xml"}, []string{".xml"}, &XMLParser{})
}

// XMLParser is a Parser implementation to handle xml files.
type XMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *XMLParser) FromBytes(byteData []byte) (dasel.Value, error) {
	if byteData == nil {
		return dasel.Value{}, fmt.Errorf("cannot parse nil xml data")
	}
	if len(byteData) == 0 || strings.TrimSpace(string(byteData)) == "" {
		return dasel.Value{}, nil
	}
	data, err := mxj.NewMapXml(byteData)
	if err != nil {
		return dasel.Value{}, fmt.Errorf("could not unmarshal data: %w", err)
	}
	return dasel.ValueOf(map[string]interface{}(data)).WithMetadata("isSingleDocument", true), nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *XMLParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)

	prettyPrint := true
	colourise := false
	indent := "  "

	for _, o := range options {
		switch o.Key {
		case OptionIndent:
			if value, ok := o.Value.(string); ok {
				indent = value
			}
		case OptionPrettyPrint:
			if value, ok := o.Value.(bool); ok {
				prettyPrint = value
			}
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				colourise = value
			}
		}
	}

	writeMap := func(val interface{}) error {
		if m, ok := val.(map[string]interface{}); ok {
			mv := mxj.New()
			for k, v := range m {
				mv[k] = v
			}

			var byteData []byte
			var err error
			if prettyPrint {
				byteData, err = mv.XmlIndent("", indent)
			} else {
				byteData, err = mv.Xml()
			}

			if err != nil {
				return err
			}
			buf.Write(byteData)
			buf.Write([]byte("\n"))
			return nil
		}
		buf.Write([]byte(fmt.Sprintf("%v\n", val)))
		return nil
	}

	switch {
	case value.Metadata("isSingleDocument") == true:
		if err := writeMap(value.Interface()); err != nil {
			return nil, err
		}
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			if err := writeMap(value.Index(i).Interface()); err != nil {
				return nil, err
			}
		}
	default:
		if err := writeMap(value.Interface()); err != nil {
			return nil, err
		}
	}

	if colourise {
		if err := ColouriseBuffer(buf, "xml"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buf.Bytes(), nil
}
