package json

import (
	"bytes"
	// "encoding/json"
	"fmt"
	json "github.com/goccy/go-json"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"io"
	"strings"
)

var _ parsing.Writer = (*jsonWriter)(nil)

// NewJSONWriter creates a new JSON writer.
func newJSONWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &jsonWriter{
		options: options,
	}, nil
}

type jsonWriter struct {
	options parsing.WriterOptions
}

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

			if _, err := fmt.Fprintf(w, `"%s": `, kv.Key); err != nil {
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
