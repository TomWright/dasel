package ini

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
	"gopkg.in/ini.v1"
)

var _ parsing.Reader = (*iniReader)(nil)

func init() {
	parsing.RegisterReader(INI, newINIReader)
	parsing.RegisterWriter(INI, newINIWriter)
}

func newINIReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &iniReader{}, nil
}

type iniReader struct{}

// Read reads a value from a byte slice.
func (j *iniReader) Read(data []byte) (*model.Value, error) {
	f, err := ini.LoadSources(ini.LoadOptions{}, data)
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

		if err := res.SetMapKey(keyName, model.NewStringValue(keyValue)); err != nil {
			return nil, err
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
