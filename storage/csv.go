package storage

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/dencoding"
	"github.com/tomwright/dasel/v2/util"
	"sort"
)

func init() {
	registerReadParser([]string{"csv"}, []string{".csv"}, &CSVParser{})
	registerWriteParser([]string{"csv"}, []string{".csv"}, &CSVParser{})
}

// CSVParser is a Parser implementation to handle csv files.
type CSVParser struct {
}

// CSVDocument represents a CSV file.
// This is required to keep headers in the expected order.
type CSVDocument struct {
	Value   []map[string]interface{}
	Headers []string
}

// FromBytes returns some data that is represented by the given bytes.
func (p *CSVParser) FromBytes(byteData []byte) (dasel.Value, error) {
	if byteData == nil {
		return dasel.Value{}, fmt.Errorf("could not read csv file: no data")
	}

	reader := csv.NewReader(bytes.NewBuffer(byteData))
	res := make([]map[string]interface{}, 0)
	records, err := reader.ReadAll()
	if err != nil {
		return dasel.Value{}, fmt.Errorf("could not read csv file: %w", err)
	}
	if len(records) == 0 {
		return dasel.Value{}, nil
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

	return dasel.ValueOf(res).
		WithMetadata("csvHeaders", headers), nil
}

func interfaceToCSVDocument(val interface{}) (*CSVDocument, error) {
	switch v := val.(type) {
	case map[string]interface{}:
		headers := make([]string, 0)
		for k := range v {
			headers = append(headers, k)
		}
		sort.Strings(headers)
		return &CSVDocument{
			Value:   []map[string]interface{}{v},
			Headers: headers,
		}, nil

	case *dencoding.Map:
		return &CSVDocument{
			Value:   []map[string]any{v.UnorderedData()},
			Headers: v.Keys(),
		}, nil

	case []interface{}:
		mapVals := make([]map[string]interface{}, 0)
		headers := make([]string, 0)
		for _, val := range v {
			switch x := val.(type) {
			case map[string]any:
				mapVals = append(mapVals, x)
				for objectKey := range x {
					found := false
					for _, existingHeader := range headers {
						if existingHeader == objectKey {
							found = true
							break
						}
					}
					if !found {
						headers = append(headers, objectKey)
					}
				}
			case *dencoding.Map:
				mapVals = append(mapVals, x.UnorderedData())
				for _, objectKey := range x.Keys() {
					found := false
					for _, existingHeader := range headers {
						if existingHeader == objectKey {
							found = true
							break
						}
					}
					if !found {
						headers = append(headers, objectKey)
					}
				}
			}
		}
		sort.Strings(headers)
		return &CSVDocument{
			Value:   mapVals,
			Headers: headers,
		}, nil

	default:
		return nil, fmt.Errorf("CSVParser.toBytes cannot handle type %T", val)
	}
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *CSVParser) ToBytes(value dasel.Value, options ...ReadWriteOption) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	// Allow for multi document output by just appending documents on the end of each other.
	// This is really only supported so as we have nicer output when converting to CSV from
	// other multi-document formats.

	docs := make([]*CSVDocument, 0)

	// headers, _ := value.Metadata("csvHeaders").([]string)

	switch {
	case value.Metadata("isSingleDocument") == true:
		doc, err := interfaceToCSVDocument(value.Interface())
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	case value.Metadata("isMultiDocument") == true:
		for i := 0; i < value.Len(); i++ {
			doc, err := interfaceToCSVDocument(value.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			docs = append(docs, doc)
		}
	default:
		doc, err := interfaceToCSVDocument(value.Interface())
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	for _, doc := range docs {
		if err := p.toBytesHandleDoc(writer, doc); err != nil {
			return nil, err
		}
	}

	return append(buffer.Bytes()), nil
}

func (p *CSVParser) toBytesHandleDoc(writer *csv.Writer, doc *CSVDocument) error {
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
				return fmt.Errorf("could not write headers: %w", err)
			}
		}

		values := make([]string, 0)
		for _, header := range doc.Headers {
			val, ok := r[header]
			if !ok {
				val = ""
			}
			values = append(values, util.ToString(val))
		}

		if err := writer.Write(values); err != nil {
			return fmt.Errorf("could not write values: %w", err)
		}

		writer.Flush()
	}

	return nil
}
