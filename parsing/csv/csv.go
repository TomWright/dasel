package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// CSV represents the CSV file format.
const CSV parsing.Format = "csv"

var _ parsing.Reader = (*csvReader)(nil)
var _ parsing.Writer = (*csvWriter)(nil)

func init() {
	parsing.RegisterReader(CSV, newCSVReader)
	parsing.RegisterWriter(CSV, newCSVWriter)
}

func newCSVReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	r := &csvReader{
		separator: ',',
	}
	if v, ok := options.Ext["csv-delimiter"]; ok && v != "" {
		r.separator = rune(v[0])
	}
	return r, nil
}

func newCSVWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	w := &csvWriter{
		separator: ',',
	}
	if v, ok := options.Ext["csv-delimiter"]; ok && v != "" {
		w.separator = rune(v[0])
	}
	return w, nil
}

type csvReader struct {
	separator rune
}

// Read reads a value from a byte slice.
func (j *csvReader) Read(data []byte) (*model.Value, error) {
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = j.separator

	res := model.NewSliceValue()

	var headers []string

	for rowI := 0; ; rowI++ {
		record, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if headers == nil {
			for _, header := range record {
				headers = append(headers, header)
			}
			continue
		}

		row := model.NewMapValue()
		for colI, field := range record {
			if colI >= len(headers) {
				return nil, fmt.Errorf("row %d has more columns than headers", rowI)
			}
			headerKey := headers[colI]

			colV, err := valueFromString(field)
			if err != nil {
				return nil, fmt.Errorf("failed reading column %q for row %d: %w", field, colI, err)
			}

			if err := row.SetMapKey(headerKey, colV); err != nil {
				return nil, fmt.Errorf("failed to set map key %q: %w", headerKey, err)
			}
		}

		if err := res.Append(row); err != nil {
			return nil, fmt.Errorf("failed to append row %d: %w", rowI, err)
		}
	}

	return res, nil
}

type csvWriter struct {
	separator rune
}

// Write writes a value to a byte slice.
func (j *csvWriter) Write(value *model.Value) ([]byte, error) {
	if !value.IsSlice() {
		return nil, fmt.Errorf("csv writer expects root output to be a slice/array, got %s", value.Type())
	}

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Comma = j.separator

	var headers []string

	if err := value.RangeSlice(func(i int, row *model.Value) error {
		if i == 0 {
			var err error
			headers, err = row.MapKeys()
			if err != nil {
				return fmt.Errorf("error getting map keys: %w", err)
			}
			if err := w.Write(headers); err != nil {
				return fmt.Errorf("error writing headers: %w", err)
			}
		}

		var values []string

		for _, headerKey := range headers {
			colV, err := row.GetMapKey(headerKey)
			if err != nil {
				return fmt.Errorf("error getting map key %q: %w", headerKey, err)
			}

			csvVal, err := valueToString(colV)
			if err != nil {
				return fmt.Errorf("error converting value to string: %w", err)
			}

			values = append(values, csvVal)
		}

		if err := w.Write(values); err != nil {
			return fmt.Errorf("error writing row: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error ranging slice: %w", err)
	}

	w.Flush()

	return buf.Bytes(), nil
}

func valueFromString(s string) (*model.Value, error) {
	return model.NewStringValue(s), nil
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
