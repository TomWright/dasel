package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestSlice(t *testing.T) {
	standardSlice := func() *model.Value {
		return model.NewValue([]any{"foo", "bar"})
	}

	modelSlice := func() *model.Value {
		res := model.NewSliceValue()
		if err := res.Append(model.NewValue("foo")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if err := res.Append(model.NewValue("bar")); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		return res
	}

	runTests := func(v func() *model.Value) func(t *testing.T) {
		return func(t *testing.T) {
			t.Run("IsSlice", func(t *testing.T) {
				v := v()
				if !v.IsSlice() {
					t.Errorf("expected value to be a slice")
				}
			})
			t.Run("GetSliceIndex", func(t *testing.T) {
				v := v()
				foo, err := v.GetSliceIndex(0)
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				got, err := foo.StringValue()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if got != "foo" {
					t.Errorf("expected foo, got %s", got)
				}
			})
			t.Run("SetSliceIndex", func(t *testing.T) {
				v := v()
				if err := v.SetSliceIndex(0, model.NewValue("baz")); err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				baz, err := v.GetSliceIndex(0)
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				got, err := baz.StringValue()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if got != "baz" {
					t.Errorf("expected baz, got %s", got)
				}
			})
			t.Run("Len", func(t *testing.T) {
				v := v()
				got, err := v.SliceLen()
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if got != 2 {
					t.Errorf("expected len of 2, got %d", got)
				}
			})
			t.Run("RangeSlice", func(t *testing.T) {
				v := v()
				var keys []int
				var vals []string
				err := v.RangeSlice(func(k int, v *model.Value) error {
					keys = append(keys, k)
					s, err := v.StringValue()
					if err != nil {
						return err
					}
					vals = append(vals, s)
					return nil
				})
				if err != nil {
					t.Errorf("unexpected error: %s", err)
					return
				}
				if len(keys) != 2 {
					t.Errorf("expected 2 keys, got %d", len(keys))
				}
				if len(vals) != 2 {
					t.Errorf("expected 2 vals, got %d", len(keys))
				}
				exp := []string{"foo", "bar"}

				for k, e := range exp {
					if keys[k] != k {
						t.Errorf("expected key %d, got %d", k, keys[k])
					}
					if vals[k] != e {
						t.Errorf("expected val %s, got %s", e, vals[k])
					}
				}
			})
			//t.Run("DeleteMapKey", func(t *testing.T) {
			//	v := v()
			//	if _, err := v.GetSliceIndex(1); err != nil {
			//		t.Errorf("unexpected error: %s", err)
			//		return
			//	}
			//	if err := v.DeleteSliceIndex(1); err != nil {
			//		t.Errorf("unexpected error: %s", err)
			//		return
			//	}
			//	_, err := v.GetSliceIndex(1)
			//	notFoundErr := &model.SliceIndexOutOfRange{}
			//	if !errors.As(err, &notFoundErr) {
			//		t.Errorf("expected index not found error, got %s", err)
			//	}
			//})
			t.Run("SliceIndexRange", func(t *testing.T) {
				t.Run("last element", func(t *testing.T) {
					v := v()
					s, err := v.SliceIndexRange(-1, -1)
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					length, err := s.SliceLen()
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					if length != 1 {
						t.Errorf("expected length of 1, got %d", length)
					}

					val, err := s.GetSliceIndex(0)
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					got, err := val.StringValue()
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					if got != "bar" {
						t.Errorf("expected bar, got %s", got)
					}
				})
				t.Run("first element", func(t *testing.T) {
					v := v()
					s, err := v.SliceIndexRange(0, 0)
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					length, err := s.SliceLen()
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					if length != 1 {
						t.Errorf("expected length of 1, got %d", length)
					}

					val, err := s.GetSliceIndex(0)
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					got, err := val.StringValue()
					if err != nil {
						t.Errorf("unexpected error: %s", err)
						return
					}
					if got != "foo" {
						t.Errorf("expected foo, got %s", got)
					}
				})
			})
		}
	}

	t.Run("standard slice", runTests(standardSlice))
	t.Run("model slice", runTests(modelSlice))
}
