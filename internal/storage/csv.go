package storage

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// CSVParser is a Parser implementation to handle yaml files.
type CSVParser struct {
}

// FromBytes returns some Data that is represented by the given bytes.
func (p *CSVParser) FromBytes(byteData []byte) (interface{}, error) {
	if byteData == nil {
		return nil, fmt.Errorf("could not read csv file: no data")
	}
	reader := csv.NewReader(bytes.NewBuffer(byteData))
	res := make([]map[string]string, 0)
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
		rowRes := make(map[string]string)
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
	return res, nil
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *CSVParser) ToBytes(value interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	rows, ok := value.([]map[string]string)
	if !ok {
		return []byte(fmt.Sprintf("%v\n", value)), nil
	}

	var headers []string
	for _, r := range rows {
		if headers == nil {
			headers = make([]string, 0)
			for k := range r {
				headers = append(headers, k)
			}
			if err := writer.Write(headers); err != nil {
				return nil, fmt.Errorf("could not write headers: %w", err)
			}
		}
		values := make([]string, 0)
		for _, v := range r {
			values = append(values, fmt.Sprint(v))
		}
		if err := writer.Write(values); err != nil {
			return nil, fmt.Errorf("could not write values: %w", err)
		}

		writer.Flush()
	}

	return append(buffer.Bytes()), nil
}
