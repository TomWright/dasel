package dasel

import (
	"reflect"
	"testing"
)

func sameSlice(x, y []interface{}) bool {
	if len(x) != len(y) {
		return false
	}

	return reflect.DeepEqual(x, y)
}

func selectTest(selector string, original interface{}, exp []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		c := NewContext(original, selector)

		values, err := c.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got := values.Interfaces()

		if !sameSlice(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func TestNewContext(t *testing.T) {
	t.Run("map propagation", func(t *testing.T) {
		orig := map[string]interface{}{
			"name": "Tom",
		}
		ctx := NewContext(orig, "property(name)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		s[0].Set(ValueOf("Frank"))

		exp := map[string]interface{}{
			"name": "Frank",
		}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})

	t.Run("nested map propagation", func(t *testing.T) {
		orig := map[string]interface{}{
			"name": map[string]interface{}{
				"first": "Tom",
				"last":  "Wright",
			},
		}
		ctx := NewContext(orig, "property(name).property(first)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		s[0].Set(ValueOf("Frank"))

		exp := map[string]interface{}{
			"name": map[string]interface{}{
				"first": "Frank",
				"last":  "Wright",
			},
		}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})

	t.Run("struct propagation", func(t *testing.T) {
		type User struct {
			Name string
		}

		orig := User{Name: "Tom"}
		ctx := NewContext(&orig, "property(Name)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		s[0].Set(ValueOf("Frank"))

		exp := User{Name: "Frank"}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})

	t.Run("nested struct propagation", func(t *testing.T) {
		type Name struct {
			First string
			Last  string
		}
		type User struct {
			Name Name
		}

		orig := User{Name: Name{First: "Tom", Last: "Wright"}}
		ctx := NewContext(&orig, "property(Name).property(First)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		s[0].Set(ValueOf("Frank"))

		exp := User{Name: Name{First: "Frank", Last: "Wright"}}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})

	t.Run("struct propagation", func(t *testing.T) {
		orig := []interface{}{
			"a", "b", "c",
		}
		ctx := NewContext(&orig, "index(0)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		s[0].Set(ValueOf("1"))

		exp := []interface{}{
			"1", "b", "c",
		}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})
}
