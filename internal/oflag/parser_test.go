package oflag_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/internal/oflag"
	"strconv"
	"testing"
)

type testCase struct {
	In  string
	Exp interface{}
	Err error
}

func runTestCase(tc testCase, fn oflag.ParseOverrideValueFn) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := fn(tc.In)
		if tc.Exp == nil && got != nil {
			t.Errorf("unexpected error: %v", got)
		} else if tc.Exp != nil && got == nil {
			t.Errorf("expected error `%v`, got none", tc.Exp)
		} else if !errors.Is(err, tc.Err) {
			t.Errorf("expected error `%v`, got `%v`", tc.Err, err)
		}
		if !cmp.Equal(tc.Exp, got) {
			t.Errorf("unexpected output\n%s", cmp.Diff(tc.Exp, got))
		}
	}
}

func TestStringParser(t *testing.T) {
	tests := []testCase{
		{In: "", Exp: "", Err: nil},
		{In: "a", Exp: "a", Err: nil},
		{In: "true", Exp: "true", Err: nil},
		{In: "123", Exp: "123", Err: nil},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, runTestCase(tc, oflag.StringParser))
	}
}

func TestBoolParser(t *testing.T) {
	tests := []testCase{
		{In: "", Exp: false, Err: oflag.ErrInvalidBool},
		{In: "asd", Exp: false, Err: oflag.ErrInvalidBool},
		{In: "yes", Exp: true, Err: nil},
		{In: "YES", Exp: true, Err: nil},
		{In: "y", Exp: true, Err: nil},
		{In: "Y", Exp: true, Err: nil},
		{In: "true", Exp: true, Err: nil},
		{In: "TRUE", Exp: true, Err: nil},
		{In: "t", Exp: true, Err: nil},
		{In: "T", Exp: true, Err: nil},
		{In: "1", Exp: true, Err: nil},
		{In: "no", Exp: false, Err: nil},
		{In: "NO", Exp: false, Err: nil},
		{In: "n", Exp: false, Err: nil},
		{In: "N", Exp: false, Err: nil},
		{In: "false", Exp: false, Err: nil},
		{In: "FALSE", Exp: false, Err: nil},
		{In: "f", Exp: false, Err: nil},
		{In: "F", Exp: false, Err: nil},
		{In: "0", Exp: false, Err: nil},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, runTestCase(tc, oflag.BoolParser))
	}
}

func TestIntParser(t *testing.T) {
	tests := []testCase{
		{In: "", Exp: 0, Err: strconv.ErrSyntax},
		{In: "asd", Exp: 0, Err: strconv.ErrSyntax},
		{In: "1.0", Exp: 0, Err: strconv.ErrSyntax},
		{In: "0", Exp: 0, Err: nil},
		{In: "12", Exp: 12, Err: nil},
		{In: "-1", Exp: -1, Err: nil},
		{In: "12345678", Exp: 12345678, Err: nil},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.In, runTestCase(tc, oflag.IntParser))
	}
}
