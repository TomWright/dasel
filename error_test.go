package dasel_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/tomwright/dasel/v2"
)

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		In  error
		Out string
	}{
		{In: dasel.ErrMissingPreviousNode, Out: "missing previous node"},
		{In: &dasel.UnknownComparisonOperatorErr{Operator: "<"}, Out: "unknown comparison operator: <"},
		{In: &dasel.InvalidIndexErr{Index: "1"}, Out: "invalid index: 1"},
		{In: &dasel.UnsupportedSelector{Selector: "..."}, Out: "selector is not supported here: ..."},
		{In: &dasel.ValueNotFound{
			Selector: ".name",
		}, Out: "no value found for selector: .name: <invalid reflect.Value>"},
		{In: &dasel.ValueNotFound{
			Selector:      ".name",
			PreviousValue: reflect.ValueOf(map[string]interface{}{}),
		}, Out: "no value found for selector: .name: map[]"},
		{In: &dasel.UnexpectedPreviousNilValue{Selector: ".name"}, Out: "previous value is nil: .name"},
		{In: &dasel.UnhandledCheckType{Value: ""}, Out: "unhandled check type: string"},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run("ErrorMessage", func(t *testing.T) {
			if exp, got := tc.Out, tc.In.Error(); exp != got {
				t.Errorf("expected %s, got %s", exp, got)
			}
		})
	}
}

func TestErrorsIs(t *testing.T) {
	type args struct {
		Err    error
		Target error
	}

	tests := []struct {
		In  args
		Out bool
	}{
		{
			In: args{
				Err:    &dasel.UnknownComparisonOperatorErr{},
				Target: &dasel.UnknownComparisonOperatorErr{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.UnknownComparisonOperatorErr{}),
				Target: &dasel.UnknownComparisonOperatorErr{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.UnknownComparisonOperatorErr{},
			},
			Out: false,
		},
		{
			In: args{
				Err:    &dasel.InvalidIndexErr{},
				Target: &dasel.InvalidIndexErr{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.InvalidIndexErr{}),
				Target: &dasel.InvalidIndexErr{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.InvalidIndexErr{},
			},
			Out: false,
		},
		{
			In: args{
				Err:    &dasel.UnsupportedSelector{},
				Target: &dasel.UnsupportedSelector{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.UnsupportedSelector{}),
				Target: &dasel.UnsupportedSelector{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.UnsupportedSelector{},
			},
			Out: false,
		},
		{
			In: args{
				Err:    &dasel.ValueNotFound{},
				Target: &dasel.ValueNotFound{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.ValueNotFound{}),
				Target: &dasel.ValueNotFound{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.ValueNotFound{},
			},
			Out: false,
		},
		{
			In: args{
				Err:    &dasel.UnexpectedPreviousNilValue{},
				Target: &dasel.UnexpectedPreviousNilValue{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.UnexpectedPreviousNilValue{}),
				Target: &dasel.UnexpectedPreviousNilValue{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.UnexpectedPreviousNilValue{},
			},
			Out: false,
		},
		{
			In: args{
				Err:    &dasel.UnhandledCheckType{},
				Target: &dasel.UnhandledCheckType{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    fmt.Errorf("some error: %w", &dasel.UnhandledCheckType{}),
				Target: &dasel.UnhandledCheckType{},
			},
			Out: true,
		},
		{
			In: args{
				Err:    errors.New("some error"),
				Target: &dasel.UnhandledCheckType{},
			},
			Out: false,
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run("ErrorMessage", func(t *testing.T) {
			if exp, got := tc.Out, errors.Is(tc.In.Err, tc.In.Target); exp != got {
				t.Errorf("expected %v, got %v", exp, got)
			}
		})
	}
}
