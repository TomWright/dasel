package storage

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// CSVParser is a Parser implementation to handle yaml files.
type CSVParser struct {
}

// CSVDocument represents a CSV file.
// This is required to keep headers in the expected order.
type CSVDocument struct {
	originalRequired
	Value   []map[string]interface{}
	Headers []string
}

// RealValue returns the real value that dasel should use when processing data.
func (d *CSVDocument) RealValue() interface{} {
	return d.Value
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *CSVParser) FromBytes(byteData []byte) (RealValue, error) {
	if byteData == nil {
		return nil, fmt.Errorf("could not read csv file: no data")
	}

	reader := csv.NewReader(bytes.NewBuffer(byteData))
	res := make([]map[string]interface{}, 0)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read csv file: %w", err)
	}
	var headers []string
	for i, row := range records {
		if i == 0 {
			headers = row
			continue
		}
		rowRes := make(map[string]interface{})
		allEmpty := true
		for index, val := range row {
			if val != "" {
				allEmpty = false
			}
			rowRes[headers[index]] = val
		}
		if !allEmpty {
			res = append(res, rowRes)
		}
	}
	return &CSVDocument{
		Value:   res,
		Headers: headers,
	}, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *CSVParser) ToBytes(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	doc, ok := value.(*CSVDocument)
	if !ok {
		return []byte(fmt.Sprintf("%v\n", value)), nil
	}

	// Iterate through the rows and detect any new headers.
	for _, r := range doc.Value {
		for k := range r {
			headerExists := false
			for _, header := range doc.Headers {
				if k == header {
					headerExists = true
					break
				}
			}
			if !headerExists {
				doc.Headers = append(doc.Headers, k)
			}
		}
	}

	// Iterate through the rows and write the output.
	for i, r := range doc.Value {
		if i == 0 {
			if err := writer.Write(doc.Headers); err != nil {
				return nil, fmt.Errorf("could not write headers: %w", err)
			}
		}

		values := make([]string, 0)
		for _, header := range doc.Headers {
			val, ok := r[header]
			if !ok {
				val = ""
			}
			values = append(values, fmt.Sprint(val))
		}

		if err := writer.Write(values); err != nil {
			return nil, fmt.Errorf("could not write values: %w", err)
		}

		writer.Flush()
	}

	return append(buffer.Bytes()), nil
}
