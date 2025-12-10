package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

func newXMLWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &xmlWriter{
		options: options,
	}, nil
}

type xmlWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *xmlWriter) Write(value *model.Value) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := xml.NewEncoder(buf)
	writer.Indent("", "  ")

	element, err := j.toElement(value)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to element: %w", err)
	}

	if err := writer.Encode(element); err != nil {
		return nil, err
	}

	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *xmlWriter) toElement(value *model.Value) (*xmlElement, error) {
	switch value.Type() {

	case model.TypeString:
		strVal, err := valueToString(value)
		return &xmlElement{
			Name:    "root",
			Content: strVal,
		}, err

	case model.TypeMap:
		kvs, err := value.MapKeyValues()
		if err != nil {
			return nil, err
		}

		el := &xmlElement{
			Name: "root",
		}

		for _, kv := range kvs {
			if strings.HasPrefix(kv.Key, "-") {
				attr := xmlAttr{
					Name: kv.Key[1:],
				}
				attr.Value, err = valueToString(kv.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to convert attribute %q to string: %w", attr.Name, err)
				}
				el.Attrs = append(el.Attrs, attr)
				continue
			}

			if kv.Key == "#text" {
				el.Content, err = valueToString(kv.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to convert content to string: %w", err)
				}
				continue
			}

			childEl, err := j.toElement(kv.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert child element %q to element: %w", kv.Key, err)
			}
			childEl.Name = kv.Key
			el.Children = append(el.Children, childEl)
		}

		return el, nil
	case model.TypeSlice:
		el := &xmlElement{
			Name: "root",
		}
		if err := value.RangeSlice(func(i int, value *model.Value) error {
			childEl, err := j.toElement(value)
			if err != nil {
				return err
			}
			childEl.Name = "item"
			el.Children = append(el.Children, childEl)

			return nil
		}); err != nil {
			return nil, err
		}
		return el, nil
	default:
		return nil, fmt.Errorf("xml writer does not support value type: %s", value.Type())
	}
}

func valueToString(v *model.Value) (string, error) {
	if v.IsNull() {
		return "", nil
	}

	switch v.Type() {
	case model.TypeString:
		stringValue, err := v.StringValue()
		if err != nil {
			return "", err
		}
		return stringValue, nil
	case model.TypeInt:
		i, err := v.IntValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case model.TypeFloat:
		i, err := v.FloatValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%g", i), nil
	case model.TypeBool:
		i, err := v.BoolValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%t", i), nil
	default:
		return "", fmt.Errorf("csv writer cannot format type %s to string", v.Type())
	}
}

func (e *xmlElement) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: e.Name}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	if len(e.Attrs) > 0 {
		for _, attr := range e.Attrs {
			if err := enc.EncodeToken(xml.Attr{
				Name:  xml.Name{Local: attr.Name},
				Value: attr.Value,
			}); err != nil {
				return err
			}
		}
	}

	if len(e.Content) > 0 {
		if err := enc.EncodeToken(xml.CharData(e.Content)); err != nil {
			return err
		}
	}

	for _, child := range e.Children {
		if err := enc.Encode(child); err != nil {
			return err
		}
	}

	return enc.EncodeToken(start.End())
}
