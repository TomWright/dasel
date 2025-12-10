package ini

import (
	"bytes"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"gopkg.in/ini.v1"
	"strings"
)

var _ parsing.Writer = (*iniWriter)(nil)

func newINIWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &iniWriter{
		options: options,
	}, nil
}

type iniWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *iniWriter) Write(value *model.Value) ([]byte, error) {
	if !value.IsMap() {
		return nil, fmt.Errorf("ini can only represent map values")
	}

	f := ini.Empty(ini.LoadOptions{
		AllowNestedValues: true,
	})

	if err := j.write(f, ini.DefaultSection, value); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if _, err := f.WriteTo(buf); err != nil {
		return nil, fmt.Errorf("failed to write ini: %w", err)
	}

	return buf.Bytes(), nil
}

func (j *iniWriter) write(f *ini.File, path string, value *model.Value) error {
	section, err := f.NewSection(path)
	if err != nil {
		return fmt.Errorf("failed to create section %s: %w", path, err)
	}

	nextSectionName := func(x string) string {
		path := strings.TrimSpace(strings.TrimPrefix(path, ini.DefaultSection))
		return strings.TrimSpace(path + " " + x)
	}

	switch value.Type() {
	case model.TypeMap:
		if err := value.RangeMap(func(s string, value *model.Value) error {
			switch {
			case value.IsScalar():
				strVal, err := valueToString(value)
				if err != nil {
					return fmt.Errorf("failed to convert value to string: %w", err)
				}
				_, err = section.NewKey(s, strVal)
				if err != nil {
					return fmt.Errorf("failed to create key %s: %w", s, err)
				}
				return nil

			case value.IsSlice():
				return fmt.Errorf("ini writer cannot represent slice values directly; consider using nested sections")

			case value.IsMap():
				if err := j.write(f, nextSectionName(s), value); err != nil {
					return err
				}
				return nil

			default:
				return fmt.Errorf("ini writer cannot represent value of type %s", value.Type())
			}
		}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("ini sections can only represent map values")
	}
	return nil
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
