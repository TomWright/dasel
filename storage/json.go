package storage

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/dencoding"
	"io"
)

func init() {
	registerReadParser([]string{"json"}, []string{".json"}, &JSONParser{})
	registerWriteParser([]string{"json"}, []string{".json"}, &JSONParser{})
}

// JSONParser is a Parser implementation to handle json files.
type JSONParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *JSONParser) FromBytes(byteData []byte) (dasel.Value, error) {
	res := make([]any, 0)

	decoder := dencoding.NewJSONDecoder(bytes.NewReader(byteData))

docLoop:
	for {
		var docData any
		if err := decoder.Decode(&docData); err != nil {
			if err == io.EOF {
				break docLoop
			}
			return dasel.Value{}, fmt.Errorf("could not unmarshal data: %w", err)
		} else {
			res = append(res, docData)
		}
	}

	switch len(res) {
	case 0:
		return dasel.Value{}, nil
	case 1:
		return dasel.ValueOf(res[0]).
			WithMetadata("isSingleDocument", true), nil
	default:
		return dasel.ValueOf(res).
			WithMetadata("isMultiDocument", true), nil
	}
}

type toBytesOptions struct {
	indent      string
	prefix      string
	prettyPrint bool
	colourise   bool
	escapeHTML  bool
}

func getToBytesOptions(options ...ReadWriteOption) toBytesOptions {
	res := toBytesOptions{
		indent:      "  ",
		prefix:      "",
		prettyPrint: true,
		colourise:   false,
		escapeHTML:  false,
	}

	for _, o := range options {
		switch o.Key {
		case OptionIndent:
			if value, ok := o.Value.(string); ok {
				res.indent = value
			}
		case OptionPrettyPrint:
			if value, ok := o.Value.(bool); ok {
				res.prettyPrint = value
			}
		case OptionColourise:
			if value, ok := o.Value.(bool); ok {
				res.colourise = value
			}
		case OptionEscapeHTML:
			if value, ok := o.Value.(bool); ok {
				res.escapeHTML = value
			}
		}
	}

	if !res.prettyPrint {
		res.indent = ""
	}

	return res
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *JSONParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	encoderOptions := make([]dencoding.JSONEncoderOption, 0)

	baseOptions := getToBytesOptions(options...)
	encoderOptions = append(encoderOptions, dencoding.JSONEscapeHTML(baseOptions.escapeHTML))
	encoderOptions = append(encoderOptions, dencoding.JSONEncodeIndent(baseOptions.prefix, baseOptions.indent))

	buffer := new(bytes.Buffer)
	encoder := dencoding.NewJSONEncoder(buffer, encoderOptions...)
	defer encoder.Close()

	switch {
	case value.Metadata("isSingleDocument") == true:
		if err := encoder.Encode(value.Interface()); err != nil {
			return nil, fmt.Errorf("could not encode single document: %w", err)
		}
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			if err := encoder.Encode(value.Index(i).Interface()); err != nil {
				return nil, fmt.Errorf("could not encode multi document [%d]: %w", i, err)
			}
		}
	default:
		if err := encoder.Encode(value.Interface()); err != nil {
			return nil, fmt.Errorf("could not encode default document type: %w", err)
		}
	}

	if baseOptions.colourise {
		if err := ColouriseBuffer(buffer, "json"); err != nil {
			return nil, fmt.Errorf("could not colourise output: %w", err)
		}
	}

	return buffer.Bytes(), nil
}
