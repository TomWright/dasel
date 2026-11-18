package json

import (
	"bytes"
	"errors"
	"fmt"
	json "github.com/goccy/go-json"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"strings"
)

// ErrJSONMaxDepthExceeded is returned when JSON nesting depth exceeds maxJSONDepth.
var ErrJSONMaxDepthExceeded = errors.New("json nesting depth exceeded")

const maxJSONDepth = 10_000

var _ parsing.Reader = (*jsonReader)(nil)

func newJSONReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &jsonReader{}, nil
}

type jsonReader struct{}

// tokenDecoder wraps a json.Decoder and keeps track of where the last token
// ended in data.
//
// goccy/go-json's Decoder.Token skips ',' and ':' wherever they appear, so on
// its own it reads `[1 2]`, `[1,,2]`, `[1,]` and `{"a" 1}` without an error.
// checkSeparator checks the bytes between tokens instead. The position is
// tracked here rather than taken from Decoder.InputOffset, because the
// decoder unescapes strings in place and its offset drifts after an escape.
type tokenDecoder struct {
	decoder *json.Decoder
	data    []byte
	pos     int
}

func (d *tokenDecoder) More() bool {
	return d.decoder.More()
}

// Token returns the next token from the decoder and moves pos past it.
func (d *tokenDecoder) Token() (json.Token, error) {
	t, err := d.decoder.Token()
	if err != nil {
		return t, err
	}
	d.pos = skipJSONToken(d.data, d.pos)
	return t, nil
}

// checkSeparator checks the input between the last token and the next one.
//
// sep is the separator that must come next, or 0 when no element has been
// read yet, in which case no separator may come next. end is the delimiter
// that may come instead of sep to close the container, or 0 if it may not.
// A separator must be followed by another value or key, not by a separator
// or the end of the container.
func (d *tokenDecoder) checkSeparator(sep byte, end byte) error {
	i := skipJSONSpace(d.data, d.pos)
	if i >= len(d.data) {
		// Let the decoder report the unexpected end of input.
		return nil
	}
	c := d.data[i]
	if sep == 0 {
		if c == ',' || c == ':' {
			return fmt.Errorf("unexpected %q at offset %d", c, i)
		}
		return nil
	}
	if end != 0 && c == end {
		return nil
	}
	if c != sep {
		if end == 0 {
			return fmt.Errorf("unexpected %q at offset %d, expected %q", c, i, sep)
		}
		return fmt.Errorf("unexpected %q at offset %d, expected %q or %q", c, i, sep, end)
	}
	i = skipJSONSpace(d.data, i+1)
	if i < len(d.data) {
		switch c := d.data[i]; c {
		case ',', ':', ']', '}':
			return fmt.Errorf("unexpected %q at offset %d after %q", c, i, sep)
		}
	}
	return nil
}

// skipJSONToken returns the index just after the token that starts at or
// after i, skipping whitespace, ',' and ':' before it as Decoder.Token does.
// The decoder has already validated the token, so this only needs to find
// where it ends.
func skipJSONToken(data []byte, i int) int {
	for i < len(data) && (isJSONSpace(data[i]) || data[i] == ',' || data[i] == ':') {
		i++
	}
	if i >= len(data) {
		return i
	}
	switch data[i] {
	case '{', '}', '[', ']':
		return i + 1
	case '"':
		for i++; i < len(data); i++ {
			switch data[i] {
			case '\\':
				i++
			case '"':
				return i + 1
			}
		}
		return i
	default:
		// A number, true, false or null.
		for i < len(data) && !isJSONSpace(data[i]) && !strings.ContainsRune(",:{}[]\"", rune(data[i])) {
			i++
		}
		return i
	}
}

// skipJSONSpace returns the index of the first byte at or after i that is not
// JSON whitespace.
func skipJSONSpace(data []byte, i int) int {
	for i < len(data) && isJSONSpace(data[i]) {
		i++
	}
	return i
}

func isJSONSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// Read reads a value from a byte slice.
// When the input contains multiple JSON values (NDJSON), they are returned
// as a branch-marked slice, mirroring the YAML multi-document behaviour.
func (j *jsonReader) Read(data []byte) (*model.Value, error) {
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	jsonDecoder.UseNumber()
	decoder := &tokenDecoder{decoder: jsonDecoder, data: data}

	var results []*model.Value

	for decoder.More() {
		v, err := j.decodeValue(decoder, 0)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}

	switch len(results) {
	case 0:
		return model.NewNullValue(), nil
	case 1:
		return results[0], nil
	default:
		slice := model.NewSliceValue()
		slice.MarkAsBranch()
		for _, v := range results {
			if err := slice.Append(v); err != nil {
				return nil, err
			}
		}
		return slice, nil
	}
}

func (j *jsonReader) decodeValue(decoder *tokenDecoder, depth int) (*model.Value, error) {
	t, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	switch t {
	case jsonOpenObject:
		res, err := j.decodeObject(decoder, depth+1)
		if err != nil {
			return nil, fmt.Errorf("could not decode object: %w", err)
		}
		return res, nil
	case jsonOpenArray:
		res, err := j.decodeArray(decoder, depth+1)
		if err != nil {
			return nil, fmt.Errorf("could not decode array: %w", err)
		}
		return res, nil
	default:
		return j.decodeToken(decoder, t)
	}
}

func (j *jsonReader) decodeObject(decoder *tokenDecoder, depth int) (*model.Value, error) {
	if depth > maxJSONDepth {
		return nil, ErrJSONMaxDepthExceeded
	}

	res := model.NewMapValue()

	var key any = nil
	first := true

	for {
		var err error
		switch {
		case first:
			err = decoder.checkSeparator(0, 0)
			first = false
		case key != nil:
			err = decoder.checkSeparator(':', 0)
		default:
			err = decoder.checkSeparator(',', '}')
		}
		if err != nil {
			return res, err
		}

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
			value, err := j.decodeArray(decoder, depth+1)
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
			value, err := j.decodeObject(decoder, depth+1)
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

func (j *jsonReader) decodeArray(decoder *tokenDecoder, depth int) (*model.Value, error) {
	if depth > maxJSONDepth {
		return nil, ErrJSONMaxDepthExceeded
	}

	res := model.NewSliceValue()
	first := true
	for {
		var err error
		if first {
			err = decoder.checkSeparator(0, 0)
			first = false
		} else {
			err = decoder.checkSeparator(',', ']')
		}
		if err != nil {
			return res, err
		}

		t, err := decoder.Token()
		if err != nil {
			// We don't expect an EOF here since we're in the middle of processing an object.
			return res, err
		}

		switch t {
		case jsonOpenArray:
			value, err := j.decodeArray(decoder, depth+1)
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
			value, err := j.decodeObject(decoder, depth+1)
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

func (j *jsonReader) decodeToken(decoder *tokenDecoder, t json.Token) (*model.Value, error) {
	switch tv := t.(type) {
	case json.Number:
		strNum := tv.String()
		// Numbers with a fraction or an exponent are floats.
		// E.g. 1.5, 1e3, 1E-2.
		if strings.ContainsAny(strNum, ".eE") {
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
