package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"io"
)

func newCSVReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	r := &csvReader{
		separator: ',',
	}
	if v, ok := options.Ext["csv-delimiter"]; ok && v != "" {
		r.separator = rune(v[0])
	}
	return r, nil
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
			headers = append(headers, record...)
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
