package storage_test

import (
	"github.com/tomwright/dasel/internal/storage"
	"testing"
)

func TestUnknownParserErr_Error(t *testing.T) {
	if exp, got := "unknown parser: x", (&storage.UnknownParserErr{Parser: "x"}).Error(); exp != got {
		t.Errorf("expected error %s, got %s", exp, got)
	}
}

func TestNewParserFromString(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "json", Out: &storage.JSONParser{}},
		{In: "yaml", Out: &storage.YAMLParser{}},
		{In: "bad", Out: nil, Err: &storage.UnknownParserErr{Parser: "bad"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewParserFromString(tc.In)
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

func TestNewParserFromFilename(t *testing.T) {
	tests := []struct {
		In  string
		Out storage.Parser
		Err error
	}{
		{In: "a.json", Out: &storage.JSONParser{}},
		{In: "a.yaml", Out: &storage.YAMLParser{}},
		{In: "a.yml", Out: &storage.YAMLParser{}},
		{In: "a.txt", Out: nil, Err: &storage.UnknownParserErr{Parser: ".txt"}},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, func(t *testing.T) {
			got, err := storage.NewParserFromFilename(tc.In)
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
