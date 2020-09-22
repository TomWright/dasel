package oflag_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/internal/oflag"
	"testing"
)

func TestOverrideFlag_Set(t *testing.T) {
	overrides := oflag.NewOverrideFlag(oflag.StringParser)
	if err := overrides.Set("a=asd"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := overrides.Set("b="); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := overrides.Set("c=a=b"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := overrides.Set("a.b.c=d.e.f"); err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if err := overrides.Set(""); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	exp := []*oflag.Override{
		{
			Path:  "a",
			Value: "asd",
		},
		{
			Path:  "b",
			Value: "",
		},
		{
			Path:  "c",
			Value: "a=b",
		},
		{
			Path:  "a.b.c",
			Value: "d.e.f",
		},
	}

	got := overrides.Overrides()

	if !cmp.Equal(exp, got) {
		t.Errorf("unexpected data:\n%s\n", cmp.Diff(exp, got))
	}
}

func TestCombine(t *testing.T) {
	a := oflag.NewOverrideFlag(nil)
	b := oflag.NewOverrideFlag(nil)
	c := oflag.NewOverrideFlag(nil)
	_ = a.Set("a=1")
	_ = b.Set("b=2")
	_ = c.Set("c=3")
	got := oflag.Combine(a, b, c)
	exp := []*oflag.Override{
		{Path: "a", Value: "1"},
		{Path: "b", Value: "2"},
		{Path: "c", Value: "3"},
	}
	if !cmp.Equal(exp, got) {
		t.Errorf("unexpected combined result:\n%s\n", cmp.Diff(exp, got))
	}
}
