package storage_test

import (
	"bytes"
	"errors"
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
	"strings"
	"testing"
)

func TestUnknownParserErr_Error(t *testing.T) {
	if exp, got := "unknown parser: x", (&storage.UnknownParserErr{Parser: "x"}).Error(); exp != got {
		t.Errorf("expected error %s, got %s", exp, got)
	}
}

func TestNewReadParserFromString(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "json", Out: &storage.JSONParser{}},
		{In: "yaml", Out: &storage.YAMLParser{}},
		{In: "yml", Out: &storage.YAMLParser{}},
		{In: "toml", Out: &storage.TOMLParser{}},
		{In: "xml", Out: &storage.XMLParser{}},
		{In: "csv", Out: &storage.CSVParser{}},
		{In: "bad", Out: nil, Err: &storage.UnknownParserErr{Parser: "bad"}},
		{In: "-", Out: nil, Err: &storage.UnknownParserErr{Parser: "-"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewReadParserFromString(tc.In)
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Out != got {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

func TestNewWriteParserFromString(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "json", Out: &storage.JSONParser{}},
		{In: "yaml", Out: &storage.YAMLParser{}},
		{In: "yml", Out: &storage.YAMLParser{}},
		{In: "toml", Out: &storage.TOMLParser{}},
		{In: "xml", Out: &storage.XMLParser{}},
		{In: "csv", Out: &storage.CSVParser{}},
		{In: "-", Out: &storage.PlainParser{}},
		{In: "plain", Out: &storage.PlainParser{}},
		{In: "bad", Out: nil, Err: &storage.UnknownParserErr{Parser: "bad"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewWriteParserFromString(tc.In)
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Out != got {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

func TestNewReadParserFromFilename(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "a.json", Out: &storage.JSONParser{}},
		{In: "a.yaml", Out: &storage.YAMLParser{}},
		{In: "a.yml", Out: &storage.YAMLParser{}},
		{In: "a.toml", Out: &storage.TOMLParser{}},
		{In: "a.xml", Out: &storage.XMLParser{}},
		{In: "a.csv", Out: &storage.CSVParser{}},
		{In: "a.txt", Out: nil, Err: &storage.UnknownParserErr{Parser: ".txt"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewReadParserFromFilename(tc.In)
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Out != got {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

func TestNewWriteParserFromFilename(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "a.json", Out: &storage.JSONParser{}},
		{In: "a.yaml", Out: &storage.YAMLParser{}},
		{In: "a.yml", Out: &storage.YAMLParser{}},
		{In: "a.toml", Out: &storage.TOMLParser{}},
		{In: "a.xml", Out: &storage.XMLParser{}},
		{In: "a.csv", Out: &storage.CSVParser{}},
		{In: "a.txt", Out: nil, Err: &storage.UnknownParserErr{Parser: ".txt"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewWriteParserFromFilename(tc.In)
			if tc.Err == nil && err != nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err == nil {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Err != nil && err != nil && err.Error() != tc.Err.Error() {
				t.Errorf("expected err %v, got %v", tc.Err, err)
				return
			}
			if tc.Out != got {
				t.Errorf("expected result %v, got %v", tc.Out, got)
			}
		})
	}
}

var jsonData = map[string]interface{}{
	"name": "Tom",
	"preferences": map[string]interface{}{
		"favouriteColour": "red",
	},
	"colours": []interface{}{"red", "green", "blue"},
	"colourCodes": []interface{}{
		map[string]interface{}{
			"name": "red",
			"rgb":  "ff0000",
		},
		map[string]interface{}{
			"name": "green",
			"rgb":  "00ff00",
		},
		map[string]interface{}{
			"name": "blue",
			"rgb":  "0000ff",
		},
	},
}

func TestLoadFromFile(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		data, err := storage.LoadFromFile("../../tests/assets/example.json", &storage.JSONParser{})
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := &storage.BasicSingleDocument{Value: jsonData}
		if !reflect.DeepEqual(exp, data) {
			t.Errorf("data does not match: exp %v, got %v", exp, data)
		}
	})
	t.Run("BaseFilePath", func(t *testing.T) {
		_, err := storage.LoadFromFile("x.json", &storage.JSONParser{})
		if err == nil || !strings.Contains(err.Error(), "could not open file") {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("ReaderErrHandled", func(t *testing.T) {
		if _, err := storage.Load(&storage.JSONParser{}, &failingReader{}); !errors.Is(err, errFailingReaderErr) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}

var errFailingParserErr = errors.New("i am meant to fail at parsing")

type failingParser struct {
}

func (fp *failingParser) FromBytes(_ []byte) (interface{}, error) {
	return nil, errFailingParserErr
}

func (fp *failingParser) ToBytes(_ interface{}) ([]byte, error) {
	return nil, errFailingParserErr
}

var errFailingWriterErr = errors.New("i am meant to fail at writing")

type failingWriter struct {
}

func (fp *failingWriter) Write(_ []byte) (int, error) {
	return 0, errFailingWriterErr
}

var errFailingReaderErr = errors.New("i am meant to fail at reading")

type failingReader struct {
}

func (fp *failingReader) Read(_ []byte) (n int, err error) {
	return 0, errFailingReaderErr
}

func TestWrite(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		if err := storage.Write(&storage.JSONParser{}, map[string]interface{}{"name": "Tom"}, nil, &buf); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if exp, got := `{
  "name": "Tom"
}
`, buf.String(); exp != got {
			t.Errorf("unexpected output:\n%s\ngot:\n%s", exp, got)
		}
	})

	t.Run("ParserErrHandled", func(t *testing.T) {
		var buf bytes.Buffer
		if err := storage.Write(&failingParser{}, map[string]interface{}{"name": "Tom"}, nil, &buf); !errors.Is(err, errFailingParserErr) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})

	t.Run("WriterErrHandled", func(t *testing.T) {
		if err := storage.Write(&storage.JSONParser{}, map[string]interface{}{"name": "Tom"}, nil, &failingWriter{}); !errors.Is(err, errFailingWriterErr) {
			t.Errorf("unexpected error: %v", err)
			return
		}
	})
}
