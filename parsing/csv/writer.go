package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
)

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
