package storage

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"github.com/clbanning/mxj/v2"
	"strings"
)

func init() {
	// Required for https://github.com/TomWright/dasel/issues/61
	mxj.XMLEscapeCharsDecoder(true)

	registerReadParser([]string{"xml"}, []string{".xml"}, &XMLParser{})
	registerWriteParser([]string{"xml"}, []string{".xml"}, &XMLParser{})
}

// XMLParser is a Parser implementation to handle yaml files.
type XMLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *XMLParser) FromBytes(byteData []byte) (interface{}, error) {
	if byteData == nil {
		return nil, fmt.Errorf("cannot parse nil xml data")
	}
	if len(byteData) == 0 || strings.TrimSpace(string(byteData)) == "" {
		return nil, nil
	}
	data, err := mxj.NewMapXml(byteData)
	if err != nil {
		return data, fmt.Errorf("could not unmarshal data: %w", err)
	}
	return &BasicSingleDocument{
		Value: map[string]interface{}(data),
	}, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *XMLParser) ToBytes(value interface{}, options ...ReadWriteOption) ([]byte, error) {
	buf := new(bytes.Buffer)

	colourise := false

	for _, o := range options {
		switch o.Key {
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
			byteData, err := mv.XmlIndent("", "  ")
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

	switch d := value.(type) {
	case SingleDocument:
		if err := writeMap(d.Document()); err != nil {
			return nil, err
		}
	case MultiDocument:
		for _, dd := range d.Documents() {
			if err := writeMap(dd); err != nil {
				return nil, err
			}
		}
	default:
		if err := writeMap(d); err != nil {
			return nil, err
		}
	}

	if colourise {
		source := buf.String()
		buf.Reset()
		if err := quick.Highlight(buf, source, "xml", ColouriseFormatter, ColouriseStyle); err != nil {
			return nil, fmt.Errorf("could not colourise json output: %w", err)
		}
	}

	return buf.Bytes(), nil
}
