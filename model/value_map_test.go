package model_test

import (
	"errors"
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestMap(t *testing.T) {
	standardMap := func() *model.Value {
		return model.NewValue(map[string]interface{}{
			"foo": "foo1",
			"bar": "bar1",
		})
	}

	dencodingMap := func() *model.Value {
		return model.NewValue(orderedmap.NewMap().
			Set("foo", "foo1").
			Set("bar", "bar1"))
	}

	modelMap := func() *model.Value {
		res := model.NewMapValue()
		if err := res.SetMapKey("foo", model.NewValue("foo1")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := res.SetMapKey("bar", model.NewValue("bar1")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		return res
	}

	runTests := func(v func() *model.Value) func(t *testing.T) {
		return func(t *testing.T) {
			t.Run("IsMap", func(t *testing.T) {
				v := v()
				if !v.IsMap() {
					t.Errorf("expected value to be a map")
				}
			})
			t.Run("GetMapKey", func(t *testing.T) {
				v := v()
				foo, err := v.GetMapKey("foo")
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				got, err := foo.StringValue()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if got != "foo1" {
					t.Errorf("expected foo1, got %s", got)
				}
			})
			t.Run("SetMapKey", func(t *testing.T) {
				v := v()
				if err := v.SetMapKey("baz", model.NewValue("baz1")); err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				baz, err := v.GetMapKey("baz")
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				got, err := baz.StringValue()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if got != "baz1" {
					t.Errorf("expected baz1, got %s", got)
				}
			})
			t.Run("MapKeys", func(t *testing.T) {
				v := v()
				keys, err := v.MapKeys()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if len(keys) != 2 {
					t.Errorf("expected 2 keys, got %d", len(keys))
				}
				exp := []string{"foo", "bar"}
				for _, k := range exp {
					var found bool
					for _, e := range keys {
						if e == k {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected key %s not found", k)
					}
				}
			})
			t.Run("RangeMap", func(t *testing.T) {
				v := v()
				var keys []string
				err := v.RangeMap(func(k string, v *model.Value) error {
					keys = append(keys, k)
					return nil
				})
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if len(keys) != 2 {
					t.Errorf("expected 2 keys, got %d", len(keys))
				}
				exp := []string{"foo", "bar"}
				for _, k := range exp {
					var found bool
					for _, e := range keys {
						if e == k {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected key %s not found", k)
					}
				}
			})
			t.Run("DeleteMapKey", func(t *testing.T) {
				v := v()
				if _, err := v.GetMapKey("foo"); err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if err := v.DeleteMapKey("foo"); err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				_, err := v.GetMapKey("foo")
				if !errors.As(err, &model.MapKeyNotFound{}) {
					t.Errorf("expected key not found error, got %s", err)
				}
			})
		}
	}

	t.Run("standard map", runTests(standardMap))
	t.Run("dencoding map", runTests(dencodingMap))
	t.Run("model map", runTests(modelMap))
}
