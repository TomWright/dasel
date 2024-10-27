package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"unicode"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// XML represents the XML file format.
	XML parsing.Format = "xml"
)

var _ parsing.Reader = (*xmlReader)(nil)
var _ parsing.Writer = (*xmlWriter)(nil)

//var _ parsing.Writer = (*xmlWriter)(nil)

func init() {
	parsing.RegisterReader(XML, newXMLReader)
	parsing.RegisterWriter(XML, newXMLWriter)
}

func newXMLReader() (parsing.Reader, error) {
	return &xmlReader{}, nil
}

// NewXMLWriter creates a new XML writer.
func newXMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &xmlWriter{
		options: options,
	}, nil
}

type xmlReader struct{}

// Read reads a value from a byte slice.
func (j *xmlReader) Read(data []byte) (*model.Value, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Strict = true

	el, err := j.parseElement(decoder, xml.StartElement{
		Name: xml.Name{
			Local: "root",
		},
	})
	if err != nil {
		return nil, err
	}

	return el.toModel()
}

type xmlAttr struct {
	Name  string
	Value string
}

type xmlElement struct {
	Name     string
	Attrs    []xmlAttr
	Children []*xmlElement
	Content  string
}

func (e *xmlElement) toModel() (*model.Value, error) {
	attrs := model.NewMapValue()
	for _, attr := range e.Attrs {
		if err := attrs.SetMapKey(attr.Name, model.NewStringValue(attr.Value)); err != nil {
			return nil, err
		}
	}
	res := model.NewMapValue()
	if err := res.SetMapKey("name", model.NewStringValue(e.Name)); err != nil {
		return nil, err
	}
	if err := res.SetMapKey("attrs", attrs); err != nil {
		return nil, err
	}

	if err := res.SetMapKey("content", model.NewStringValue(e.Content)); err != nil {
		return nil, err
	}
	children := model.NewSliceValue()
	for _, child := range e.Children {
		childModel, err := child.toModel()
		if err != nil {
			return nil, err
		}
		if err := children.Append(childModel); err != nil {
			return nil, err
		}
	}
	if err := res.SetMapKey("children", children); err != nil {
		return nil, err
	}
	return res, nil
}

func (j *xmlReader) parseElement(decoder *xml.Decoder, element xml.StartElement) (*xmlElement, error) {
	el := &xmlElement{
		Name:     element.Name.Local,
		Attrs:    make([]xmlAttr, 0),
		Children: make([]*xmlElement, 0),
	}

	for _, attr := range element.Attr {
		el.Attrs = append(el.Attrs, xmlAttr{
			Name:  attr.Name.Local,
			Value: attr.Value,
		})
	}

	for {
		t, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			if el.Name == "root" {
				return el, nil
			}
			return nil, fmt.Errorf("unexpected EOF")
		}

		switch t := t.(type) {
		case xml.StartElement:
			child, err := j.parseElement(decoder, t)
			if err != nil {
				return nil, err
			}
			el.Children = append(el.Children, child)
		case xml.CharData:
			if unicode.IsSpace([]rune(string(t))[0]) {
				continue
			}
			el.Content += string(t)
		case xml.EndElement:
			return el, nil
		default:
			return nil, fmt.Errorf("unexpected token: %v", t)
		}
	}
}

type xmlWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *xmlWriter) Write(value *model.Value) ([]byte, error) {
	return nil, nil
}
