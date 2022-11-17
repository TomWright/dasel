package query

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

		s, err := c.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got := s.output.Interfaces()

		if !sameSlice(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func TestNewContext(t *testing.T) {

	original := map[string]interface{}{
		"name": map[string]interface{}{
			"first": "Tom",
			"last":  "Wright",
		},
		"colours": []interface{}{
			"red", "green", "blue",
		},
	}

	t.Run(
		"SingleLevelProperty",
		selectTest(
			"name",
			original,
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"SingleLevelPropertyFunc",
		selectTest(
			"property(name)",
			original,
			[]interface{}{
				map[string]interface{}{
					"first": "Tom",
					"last":  "Wright",
				},
			},
		),
	)

	t.Run(
		"NestedPropertyFunc",
		selectTest(
			"property(name).property(first)",
			original,
			[]interface{}{
				"Tom",
			},
		),
	)

	t.Run(
		"NestedMultiPropertyFunc",
		selectTest(
			"property(name).property(first,last)",
			original,
			[]interface{}{
				"Tom",
				"Wright",
			},
		),
	)

	t.Run(
		"NestedMultiMissingPropertyFunc",
		selectTest(
			"property(name).property(first,last,middle?)",
			original,
			[]interface{}{
				"Tom",
				"Wright",
			},
		),
	)

	t.Run(
		"Index",
		selectTest(
			"colours.index(1)",
			original,
			[]interface{}{
				"green",
			},
		),
	)

	t.Run(
		"IndexShorthand",
		selectTest(
			"colours.[1]",
			original,
			[]interface{}{
				"green",
			},
		),
	)

	t.Run(
		"First",
		selectTest(
			"colours.first()",
			original,
			[]interface{}{
				"red",
			},
		),
	)

	t.Run(
		"Last",
		selectTest(
			"colours.last()",
			original,
			[]interface{}{
				"blue",
			},
		),
	)

	t.Run(
		"All",
		selectTest(
			"colours.all()",
			original,
			[]interface{}{
				"red",
				"green",
				"blue",
			},
		),
	)

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

		s.output[0].Set(ValueOf("Frank"))

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

		s.output[0].Set(ValueOf("Frank"))

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

		s.output[0].Set(ValueOf("Frank"))

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

		s.output[0].Set(ValueOf("Frank"))

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

		s.output[0].Set(ValueOf("1"))

		exp := []interface{}{
			"1", "b", "c",
		}
		got := orig

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})

	t.Run("metadata", func(t *testing.T) {
		orig := []interface{}{
			"abc", true, false, 1, 1.1, []interface{}{1},
		}
		ctx := NewContext(&orig, "all().metadata(type)")
		s, err := ctx.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		exp := []interface{}{
			"string", "bool", "bool", "int", "float64", "slice",
		}
		got := s.output.Interfaces()

		if !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	})
}
