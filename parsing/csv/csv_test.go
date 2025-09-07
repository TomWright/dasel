package csv

import (
	"github.com/tomwright/dasel/v3/model"
	"testing"
)

func TestValueToString(t *testing.T) {
	tests := []struct {
		desc string
		in   func() (*model.Value, error)
		exp  string
	}{
		{
			desc: "basic string",
			in: func() (*model.Value, error) {
				return model.NewStringValue("hello"), nil
			},
			exp: "hello",
		},
		{
			desc: "string",
			in: func() (*model.Value, error) {
				return model.NewStringValue("hello, there!!"), nil
			},
			exp: "hello, there!!",
		},
		{
			desc: "int",
			in: func() (*model.Value, error) {
				return model.NewIntValue(123), nil
			},
			exp: "123",
		},
		{
			desc: "float",
			in: func() (*model.Value, error) {
				return model.NewFloatValue(1.234), nil
			},
			exp: "1.234",
		},
		{
			desc: "null",
			in: func() (*model.Value, error) {
				return model.NewNullValue(), nil
			},
			exp: "",
		},
		{
			desc: "bool true",
			in: func() (*model.Value, error) {
				return model.NewBoolValue(true), nil
			},
			exp: "true",
		},
		{
			desc: "bool false",
			in: func() (*model.Value, error) {
				return model.NewBoolValue(false), nil
			},
			exp: "false",
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			in, err := tc.in()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := valueToString(in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.exp {
				t.Errorf("expected '%s' but got '%s'", tc.exp, got)
			}
		})
	}
}
