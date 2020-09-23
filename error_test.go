package dasel_test

import (
	"github.com/tomwright/dasel"
	"testing"
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
		{In: &dasel.UnsupportedTypeForSelector{
			Value: map[string]interface{}{},
			Selector: dasel.Selector{
				Raw:       ".a.b.c",
				Current:   ".a",
				Remaining: ".b.c",
				Type:      "INDEX",
				Index:     1,
			},
		}, Out: "selector [INDEX] does not support value: map[string]interface {}: map[]"},
		{In: &dasel.ValueNotFound{
			Selector: ".name",
			Node: &dasel.Node{
				Selector: dasel.Selector{
					Current: ".name",
				},
			},
		}, Out: "no value found for selector: .name: <nil>"},
		{In: &dasel.ValueNotFound{
			Selector: ".name",
			Node: &dasel.Node{
				Previous: &dasel.Node{
					Value: map[string]interface{}{},
				},
				Selector: dasel.Selector{
					Current: ".name",
				},
			},
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
