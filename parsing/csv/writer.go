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

	// Rows may have different keys, so the header is every key in first-seen order.
	var headers []string
	seen := make(map[string]struct{})

	if err := value.RangeSlice(func(_ int, row *model.Value) error {
		keys, err := row.MapKeys()
		if err != nil {
			return fmt.Errorf("error getting map keys: %w", err)
		}
		for _, key := range keys {
			if _, ok := seen[key]; !ok {
				seen[key] = struct{}{}
				headers = append(headers, key)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error ranging slice: %w", err)
	}

	if err := value.RangeSlice(func(i int, row *model.Value) error {
		if i == 0 {
			if err := w.Write(headers); err != nil {
				return fmt.Errorf("error writing headers: %w", err)
			}
		}

		var values []string

		for _, headerKey := range headers {
			exists, err := row.MapKeyExists(headerKey)
			if err != nil {
				return fmt.Errorf("error checking map key %q: %w", headerKey, err)
			}
			if !exists {
				values = append(values, "")
				continue
			}

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
