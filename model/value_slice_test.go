package model_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/ptr"
)

func TestNewSliceValue(t *testing.T) {
	x := model.NewSliceValue()
	if err := x.Append(model.NewStringValue("hello")); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := x.Append(model.NewStringValue("world")); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	got, err := x.SliceValue()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	exp := []any{ptr.To("hello"), ptr.To("world")}
	if !cmp.Equal(exp, got) {
		t.Errorf("unexpected result: %s", cmp.Diff(exp, got))
	}
}
