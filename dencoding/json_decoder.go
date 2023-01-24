package dencoding

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

// JSONDecoder wraps a standard json encoder to implement custom ordering logic.
type JSONDecoder struct {
	decoder *json.Decoder
}

// NewJSONDecoder returns a new dencoding JSONDecoder.
func NewJSONDecoder(r io.Reader, options ...JSONDecoderOption) *JSONDecoder {
	jsonDecoder := json.NewDecoder(r)
	decoder := &JSONDecoder{
		decoder: jsonDecoder,
	}
	for _, o := range options {
		o.ApplyDecoder(decoder)
	}
	return decoder
}

// Decode decodes the next item found in the decoder and writes it to v.
func (decoder *JSONDecoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("invalid decode target: %s", reflect.TypeOf(v))
	}

	rve := rv.Elem()

	t, err := decoder.decoder.Token()
	if err != nil {
		return err
	}

	switch t {
	case jsonOpenObject:
		object, err := decoder.decodeObject()
		if err != nil {
			return fmt.Errorf("could not decode object: %w", err)
		}
		rve.Set(reflect.ValueOf(object))
	case jsonOpenArray:
		arr, err := decoder.decodeArray()
		if err != nil {
			return fmt.Errorf("could not decode array: %w", err)
		}
		rve.Set(reflect.ValueOf(arr))
	default:
		rve.Set(reflect.ValueOf(t))
	}

	return nil
}

func (decoder *JSONDecoder) decodeObject() (*Map, error) {
	res := NewMap()

	var key any = nil

	for {
		t, err := decoder.decoder.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case jsonOpenArray:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := decoder.decodeArray()
			if err != nil {
				return res, err
			}
			res.Set(key.(string), value)
			key = nil
		case jsonCloseArray:
			return res, fmt.Errorf("unexpected token: %v", t)
		case jsonCloseObject:
			return res, nil
		case jsonOpenObject:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := decoder.decodeObject()
			if err != nil {
				return res, err
			}
			res.Set(key.(string), value)
			key = nil
		default:
			if key == nil {
				key = t
			} else {
				res.Set(key.(string), t)
				key = nil
			}
		}
	}
}

func (decoder *JSONDecoder) decodeArray() ([]any, error) {
	res := make([]any, 0)
	for {
		t, err := decoder.decoder.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case jsonOpenArray:
			value, err := decoder.decodeArray()
			if err != nil {
				return res, err
			}
			res = append(res, value)
		case jsonCloseArray:
			return res, nil
		case jsonCloseObject:
			return res, fmt.Errorf("unexpected token: %t", t)
		case jsonOpenObject:
			value, err := decoder.decodeObject()
			if err != nil {
				return res, err
			}
			res = append(res, value)
		default:
			res = append(res, t)
		}
	}
}
