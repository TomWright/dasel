package ini

import (
	"bytes"
	"fmt"
	"gopkg.in/ini.v1"
	"strings"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// INI represents the ini file format.
	INI parsing.Format = "ini"
)

var _ parsing.Reader = (*iniReader)(nil)
var _ parsing.Writer = (*iniWriter)(nil)

func init() {
	parsing.RegisterReader(INI, newINIReader)
	parsing.RegisterWriter(INI, newINIWriter)
}

func newINIReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &iniReader{}, nil
}

// NewINIWriter creates a new INI writer.
func newINIWriter(options parsing.WriterOptions) (parsing.Writer, error) {
	return &iniWriter{
		options: options,
	}, nil
}

type iniReader struct{}

// Read reads a value from a byte slice.
func (j *iniReader) Read(data []byte) (*model.Value, error) {
	f, err := ini.Load(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ini: %w", err)
	}

	res, err := j.readSection(f.Section(ini.DefaultSection))
	if err != nil {
		return nil, err
	}

	for _, s := range f.Sections() {
		if s.Name() == ini.DefaultSection {
			continue
		}
		sectionValue, err := j.readSection(s)
		if err != nil {
			return nil, err
		}
		if err := res.SetMapKey(s.Name(), sectionValue); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (j *iniReader) readSection(s *ini.Section) (*model.Value, error) {
	res := model.NewMapValue()
	for _, k := range s.Keys() {
		keyName := k.Name()
		keyValue := k.Value()

		if strings.HasSuffix(keyName, "[]") {
			keyName = strings.TrimSuffix(keyName, "[]")
			keyExists, err := res.MapKeyExists(keyName)
			if err != nil {
				return nil, err
			}
			var sliceValue *model.Value
			if keyExists {
				sliceValue, err = res.GetMapKey(keyName)
				if err != nil {
					return nil, err
				}
			} else {
				sliceValue = model.NewSliceValue()
				if err := res.SetMapKey(keyName, sliceValue); err != nil {
					return nil, err
				}
			}
			if err := sliceValue.Append(model.NewStringValue(keyValue)); err != nil {
				return nil, err
			}
		} else {
			if err := res.SetMapKey(keyName, model.NewStringValue(keyValue)); err != nil {
				return nil, err
			}
		}
	}
	for _, s := range s.ChildSections() {
		childSection, err := j.readSection(s)
		if err != nil {
			return nil, err
		}
		if err := res.SetMapKey(s.Name(), childSection); err != nil {
			return nil, err
		}
	}
	return res, nil
}

type iniWriter struct {
	options parsing.WriterOptions
}

// Write writes a value to a byte slice.
func (j *iniWriter) Write(value *model.Value) ([]byte, error) {
	if !value.IsMap() {
		return nil, fmt.Errorf("ini can only represent map values")
	}

	f := ini.Empty(ini.LoadOptions{})

	defaultSection, err := f.NewSection(ini.DefaultSection)
	if err != nil {
		return nil, fmt.Errorf("failed to create default section: %w", err)
	}

	if err := j.writeSection(defaultSection, value); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	if _, err := f.WriteTo(buf); err != nil {
		return nil, fmt.Errorf("failed to write ini: %w", err)
	}

	return buf.Bytes(), nil
}

func (j *iniWriter) writeSection(section *ini.Section, value *model.Value) error {
	switch value.Type() {
	case model.TypeMap:
		if err := value.RangeMap(func(s string, value *model.Value) error {
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
