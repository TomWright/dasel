package csv

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"unicode/utf8"
)

// CSV represents the CSV file format.
const CSV parsing.Format = "csv"

var _ parsing.Reader = (*csvReader)(nil)
var _ parsing.Writer = (*csvWriter)(nil)

func init() {
	parsing.RegisterReader(CSV, newCSVReader)
	parsing.RegisterWriter(CSV, newCSVWriter)
}

// separatorFromOption decodes the csv-delimiter option as a rune. Taking
// v[0] gives the first byte, so a multi-byte delimiter such as § became the
// first half of its encoding and silently matched nothing.
func separatorFromOption(ext map[string]string, def rune) (rune, error) {
	v, ok := ext["csv-delimiter"]
	if !ok || v == "" {
		return def, nil
	}
	r, size := utf8.DecodeRuneInString(v)
	if r == utf8.RuneError || size != len(v) {
		return 0, fmt.Errorf("csv-delimiter must be a single character, got %q", v)
	}
	return r, nil
}

func newCSVWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	separator, err := separatorFromOption(options.Ext, ',')
	if err != nil {
		return nil, err
	}
	return &csvWriter{separator: separator}, nil
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
