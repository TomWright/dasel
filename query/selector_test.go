package query

import (
	"reflect"
	"testing"
)

func collectAll(r SelectorResolver) ([]Selector, error) {
	res := make([]Selector, 0)

	for {
		s, err := r.Next()
		if err != nil {
			return res, err
		}
		if s == nil {
			break
		}
		res = append(res, *s)
	}

	return res, nil
}

func TestStandardSelectorResolver_Next(t *testing.T) {
	r := NewSelectorResolver("index(1).property(user).name.property(first,last?)", nil)

	got, err := collectAll(r)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	exp := []Selector{
		{
			funcName: "index",
			funcArgs: []string{"1"},
		},
		{
			funcName: "property",
			funcArgs: []string{"user"},
		},
		{
			funcName: "property",
			funcArgs: []string{"name"},
		},
		{
			funcName: "property",
			funcArgs: []string{"first", "last?"},
		},
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("exp: %v, got: %v", exp, got)
	}
}

func TestStandardSelectorResolver_Next_Nested(t *testing.T) {
	r := NewSelectorResolver("nested(a().b(),c(),d()).nested(a().b(),c(),d())", nil)

	got, err := collectAll(r)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	exp := []Selector{
		{
			funcName: "nested",
			funcArgs: []string{"a().b()", "c()", "d()"},
		},
		{
			funcName: "nested",
			funcArgs: []string{"a().b()", "c()", "d()"},
		},
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("exp: %v, got: %v", exp, got)
	}
}
