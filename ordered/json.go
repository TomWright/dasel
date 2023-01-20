package ordered

import (
	"encoding/json"
	"fmt"
)

const (
	openObject  = json.Delim('{')
	closeObject = json.Delim('}')
	openArray   = json.Delim('[')
	closeArray  = json.Delim(']')
)

// MarshalJSON marshals the Map into JSON bytes.
func (m *Map) MarshalJSON() ([]byte, error) {
	// The standard json package sorts map keys during the marshal process.
	// This is frustrating. To get around it, a custom JSON encoder will be needed.
	return json.Marshal(m.data)
}

func UnmarshalJSON(dec *json.Decoder) (any, error) {
	for {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		switch t {
		case openObject:
			return UnmarshalJSONObject(dec)
		case openArray:
			return UnmarshalJSONArray(dec)
		default:
			return t, nil
		}
	}
}

func UnmarshalJSONObject(dec *json.Decoder) (*Map, error) {
	res := NewMap()

	var key any = nil

	for {
		t, err := dec.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case openArray:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := UnmarshalJSONArray(dec)
			if err != nil {
				return res, err
			}
			res.Set(key.(string), value)
			key = nil
		case closeArray:
			return res, fmt.Errorf("unexpected token: %v", t)
		case closeObject:
			return res, nil
		case openObject:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := UnmarshalJSONObject(dec)
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

func UnmarshalJSONArray(dec *json.Decoder) ([]any, error) {
	res := make([]any, 0)
	for {
		t, err := dec.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case openArray:
			value, err := UnmarshalJSONArray(dec)
			if err != nil {
				return res, err
			}
			res = append(res, value)
		case closeArray:
			return res, nil
		case closeObject:
			return res, fmt.Errorf("unexpected token: %t", t)
		case openObject:
			value, err := UnmarshalJSONObject(dec)
			if err != nil {
				return res, err
			}
			res = append(res, value)
		default:
			res = append(res, t)
		}
	}
}
