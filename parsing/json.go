package parsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

const (
	jsonOpenObject  = json.Delim('{')
	jsonCloseObject = json.Delim('}')
	jsonOpenArray   = json.Delim('[')
	jsonCloseArray  = json.Delim(']')
)

// NewJSONReader creates a new JSON reader.
func NewJSONReader() (Reader, error) {
	return &jsonReader{}, nil
}

// NewJSONWriter creates a new JSON writer.
func NewJSONWriter() (Writer, error) {
	return &jsonWriter{}, nil
}

type jsonReader struct{}

// Read reads a value from a byte slice.
func (j *jsonReader) Read(data []byte) (*model.Value, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	t, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	var res *model.Value

	switch t {
	case jsonOpenObject:
		res, err = j.decodeObject(decoder)
		if err != nil {
			return nil, fmt.Errorf("could not decode object: %w", err)
		}
	case jsonOpenArray:
		res, err = j.decodeArray(decoder)
		if err != nil {
			return nil, fmt.Errorf("could not decode array: %w", err)
		}
	default:
		res, err = j.decodeToken(decoder, t)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (j *jsonReader) decodeObject(decoder *json.Decoder) (*model.Value, error) {
	res := model.NewMapValue()

	var key any = nil

	for {
		t, err := decoder.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case jsonOpenArray:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := j.decodeArray(decoder)
			if err != nil {
				return res, err
			}
			if err := res.SetMapKey(key.(string), value); err != nil {
				return res, err
			}
			key = nil
		case jsonCloseArray:
			return res, fmt.Errorf("unexpected token: %v", t)
		case jsonCloseObject:
			return res, nil
		case jsonOpenObject:
			if key == nil {
				return res, fmt.Errorf("unexpected token: %v", t)
			}
			value, err := j.decodeObject(decoder)
			if err != nil {
				return res, err
			}
			if err := res.SetMapKey(key.(string), value); err != nil {
				return res, err
			}
			key = nil
		default:
			if key == nil {
				if tStr, ok := t.(string); ok {
					key = tStr
				} else {
					return nil, fmt.Errorf("unexpected token: %v", t)
				}
			} else {
				value, err := j.decodeToken(decoder, t)
				if err != nil {
					return nil, err
				}
				if err := res.SetMapKey(key.(string), value); err != nil {
					return res, err
				}
				key = nil
			}
		}
	}
}

func (j *jsonReader) decodeArray(decoder *json.Decoder) (*model.Value, error) {
	res := model.NewSliceValue()
	for {
		t, err := decoder.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case jsonOpenArray:
			value, err := j.decodeArray(decoder)
			if err != nil {
				return res, err
			}
			if err := res.Append(value); err != nil {
				return res, err
			}
		case jsonCloseArray:
			return res, nil
		case jsonCloseObject:
			return res, fmt.Errorf("unexpected token: %t", t)
		case jsonOpenObject:
			value, err := j.decodeObject(decoder)
			if err != nil {
				return res, err
			}
			if err := res.Append(value); err != nil {
				return res, err
			}
		default:
			value, err := j.decodeToken(decoder, t)
			if err != nil {
				return nil, err
			}
			if err := res.Append(value); err != nil {
				return res, err
			}
		}
	}
}

func (j *jsonReader) decodeToken(decoder *json.Decoder, t json.Token) (*model.Value, error) {
	switch tv := t.(type) {
	case json.Number:
		strNum := tv.String()
		if strings.Contains(strNum, ".") {
			floatNum, err := tv.Float64()
			if err == nil {
				return model.NewFloatValue(floatNum), nil
			}
			return nil, err
		}
		intNum, err := tv.Int64()
		if err == nil {
			return model.NewIntValue(intNum), nil
		}

		return nil, err
	default:
		return model.NewValue(tv), nil
	}
}

type jsonWriter struct{}

// Write writes a value to a byte slice.
func (j *jsonWriter) Write(value *model.Value) ([]byte, error) {
	buf := new(bytes.Buffer)

	es := encoderState{indentStr: "    "}

	encoderFn := func(v any) error {
		res, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = buf.Write(res)
		return err
	}

	if err := j.write(buf, encoderFn, es, value); err != nil {
		return nil, err
	}

	if _, err := buf.Write([]byte("\n")); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type encoderState struct {
	indent    int
	indentStr string
}

func (es encoderState) inc() encoderState {
	es.indent++
	return es
}

func (es encoderState) writeIndent(w io.Writer) error {
	if es.indent == 0 || es.indentStr == "" {
		return nil
	}
	i := strings.Repeat(es.indentStr, es.indent)
	if _, err := w.Write([]byte(i)); err != nil {
		return err
	}
	return nil
}

type encoderFn func(v any) error

func (j *jsonWriter) write(w io.Writer, encoder encoderFn, es encoderState, value *model.Value) error {
	switch value.Type() {
	case model.TypeMap:
		return j.writeMap(w, encoder, es, value)
	case model.TypeSlice:
		return j.writeSlice(w, encoder, es, value)
	case model.TypeString:
		val, err := value.StringValue()
		if err != nil {
			return err
		}
		return encoder(val)
	case model.TypeInt:
		val, err := value.IntValue()
		if err != nil {
			return err
		}
		return encoder(val)
	case model.TypeFloat:
		val, err := value.FloatValue()
		if err != nil {
			return err
		}
		return encoder(val)
	case model.TypeBool:
		val, err := value.BoolValue()
		if err != nil {
			return err
		}
		return encoder(val)
	case model.TypeNull:
		return encoder(nil)
	default:
		return fmt.Errorf("unsupported type: %s", value.Type())
	}
}

func (j *jsonWriter) writeMap(w io.Writer, encoder encoderFn, es encoderState, value *model.Value) error {
	kvs, err := value.MapKeyValues()
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(`{`)); err != nil {
		return err
	}

	if len(kvs) > 0 {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}

		incEs := es.inc()
		for i, kv := range kvs {
			if err := incEs.writeIndent(w); err != nil {
				return err
			}

			if _, err := w.Write([]byte(fmt.Sprintf(`"%s": `, kv.Key))); err != nil {
				return err
			}

			if err := j.write(w, encoder, incEs, kv.Value); err != nil {
				return err
			}

			if i < len(kvs)-1 {
				if _, err := w.Write([]byte(`,`)); err != nil {
					return err
				}
			}

			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		if err := es.writeIndent(w); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(`}`)); err != nil {
		return err
	}

	return nil
}

func (j *jsonWriter) writeSlice(w io.Writer, encoder encoderFn, es encoderState, value *model.Value) error {
	if _, err := w.Write([]byte(`[`)); err != nil {
		return err
	}

	length, err := value.SliceLen()
	if err != nil {
		return fmt.Errorf("error getting slice length: %w", err)
	}

	if length > 0 {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
		incEs := es.inc()
		for i := 0; i < length; i++ {
			if err := incEs.writeIndent(w); err != nil {
				return err
			}
			va, err := value.GetSliceIndex(i)
			if err != nil {
				return fmt.Errorf("error getting slice index %d: %w", i, err)
			}
			if err := j.write(w, encoder, incEs, va); err != nil {
				return err
			}
			if i < length-1 {
				if _, err := w.Write([]byte(`,`)); err != nil {
					return err
				}
			}
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		if err := es.writeIndent(w); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(`]`)); err != nil {
		return err
	}

	return nil
}
