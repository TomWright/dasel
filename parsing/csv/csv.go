package csv

import (
	"fmt"
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

func newCSVWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	w := &csvWriter{
		separator: ',',
	}
	if v, ok := options.Ext["csv-delimiter"]; ok && v != "" {
		w.separator = rune(v[0])
	}
	return w, nil
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
