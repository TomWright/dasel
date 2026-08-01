package plain

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

var _ parsing.Writer = (*plainWriter)(nil)

func newPlainWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &plainWriter{options: options}, nil
}

type plainWriter struct {
	options parsing.WriterOptions
}

// Separator returns an empty document separator. Each document already ends in a
// newline, so multi-document output is one value per line with no blank lines.
func (w *plainWriter) Separator() []byte {
	return nil
}

// Write writes a scalar value as an unquoted string.
// Maps and slices have no plain representation and are rejected.
func (w *plainWriter) Write(value *model.Value) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("nil value")
	}

	str, err := scalarToString(value)
	if err != nil {
		return nil, err
	}

	return []byte(str + "\n"), nil
}

// scalarToString renders a scalar value as an unquoted string, matching the
// conversions used by the toString function.
func scalarToString(value *model.Value) (string, error) {
	switch value.Type() {
	case model.TypeString:
		return value.StringValue()
	case model.TypeInt:
		i, err := value.IntValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", i), nil
	case model.TypeFloat:
		f, err := value.FloatValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%g", f), nil
	case model.TypeBool:
		b, err := value.BoolValue()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", b), nil
	case model.TypeNull:
		// Null has no text form. Write an empty line rather than the literal
		// "null", which would be indistinguishable from the string "null".
		return "", nil
	default:
		return "", fmt.Errorf("plain can only represent scalar values, got %s", value.Type())
	}
}
