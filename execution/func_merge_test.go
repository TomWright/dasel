package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncMerge(t *testing.T) {
	t.Run("flat maps", testCase{
		inFn: func() *model.Value {
			a := model.NewMapValue()
			if err := a.SetMapKey("foo", model.NewStringValue("afoo")); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := a.SetMapKey("bar", model.NewStringValue("abar")); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			b := model.NewMapValue()
			if err := b.SetMapKey("bar", model.NewStringValue("bbar")); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := b.SetMapKey("baz", model.NewStringValue("bbaz")); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			res := model.NewMapValue()
			if err := res.SetMapKey("a", a); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if err := res.SetMapKey("b", b); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			return res
		},
		s: `merge(a, b)`,
		outFn: func() *model.Value {
			b := model.NewMapValue()
			if err := b.SetMapKey("foo", model.NewStringValue("afoo")); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := b.SetMapKey("bar", model.NewStringValue("bbar")); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := b.SetMapKey("baz", model.NewStringValue("bbaz")); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return b
		},
	}.run)

	t.Run("deep nested maps", testCase{
		inFn: func() *model.Value {
			innerA := model.NewMapValue()
			_ = innerA.SetMapKey("x", model.NewIntValue(1))
			_ = innerA.SetMapKey("y", model.NewIntValue(2))
			a := model.NewMapValue()
			_ = a.SetMapKey("nested", innerA)

			innerB := model.NewMapValue()
			_ = innerB.SetMapKey("y", model.NewIntValue(3))
			_ = innerB.SetMapKey("z", model.NewIntValue(4))
			b := model.NewMapValue()
			_ = b.SetMapKey("nested", innerB)

			res := model.NewMapValue()
			_ = res.SetMapKey("a", a)
			_ = res.SetMapKey("b", b)
			return res
		},
		s: `merge(a, b)`,
		outFn: func() *model.Value {
			inner := model.NewMapValue()
			_ = inner.SetMapKey("x", model.NewIntValue(1))
			_ = inner.SetMapKey("y", model.NewIntValue(3))
			_ = inner.SetMapKey("z", model.NewIntValue(4))
			out := model.NewMapValue()
			_ = out.SetMapKey("nested", inner)
			return out
		},
	}.run)

	t.Run("deep three levels", testCase{
		inFn: func() *model.Value {
			l2a := model.NewMapValue()
			_ = l2a.SetMapKey("deep", model.NewStringValue("a"))
			l1a := model.NewMapValue()
			_ = l1a.SetMapKey("l2", l2a)
			a := model.NewMapValue()
			_ = a.SetMapKey("l1", l1a)

			l2b := model.NewMapValue()
			_ = l2b.SetMapKey("deep", model.NewStringValue("b"))
			_ = l2b.SetMapKey("extra", model.NewStringValue("val"))
			l1b := model.NewMapValue()
			_ = l1b.SetMapKey("l2", l2b)
			b := model.NewMapValue()
			_ = b.SetMapKey("l1", l1b)

			res := model.NewMapValue()
			_ = res.SetMapKey("a", a)
			_ = res.SetMapKey("b", b)
			return res
		},
		s: `merge(a, b)`,
		outFn: func() *model.Value {
			l2 := model.NewMapValue()
			_ = l2.SetMapKey("deep", model.NewStringValue("b"))
			_ = l2.SetMapKey("extra", model.NewStringValue("val"))
			l1 := model.NewMapValue()
			_ = l1.SetMapKey("l2", l2)
			out := model.NewMapValue()
			_ = out.SetMapKey("l1", l1)
			return out
		},
	}.run)

	t.Run("deep with slice replacement", testCase{
		inFn: func() *model.Value {
			innerA := model.NewMapValue()
			sliceA := model.NewSliceValue()
			_ = sliceA.Append(model.NewIntValue(1))
			_ = sliceA.Append(model.NewIntValue(2))
			_ = innerA.SetMapKey("items", sliceA)
			a := model.NewMapValue()
			_ = a.SetMapKey("nested", innerA)

			innerB := model.NewMapValue()
			sliceB := model.NewSliceValue()
			_ = sliceB.Append(model.NewIntValue(3))
			_ = innerB.SetMapKey("items", sliceB)
			b := model.NewMapValue()
			_ = b.SetMapKey("nested", innerB)

			res := model.NewMapValue()
			_ = res.SetMapKey("a", a)
			_ = res.SetMapKey("b", b)
			return res
		},
		s: `merge(a, b)`,
		outFn: func() *model.Value {
			inner := model.NewMapValue()
			slice := model.NewSliceValue()
			_ = slice.Append(model.NewIntValue(3))
			_ = inner.SetMapKey("items", slice)
			out := model.NewMapValue()
			_ = out.SetMapKey("nested", inner)
			return out
		},
	}.run)

	t.Run("deep mixed types", testCase{
		inFn: func() *model.Value {
			innerA := model.NewMapValue()
			_ = innerA.SetMapKey("x", model.NewIntValue(1))
			a := model.NewMapValue()
			_ = a.SetMapKey("key", innerA)

			b := model.NewMapValue()
			_ = b.SetMapKey("key", model.NewStringValue("scalar"))

			res := model.NewMapValue()
			_ = res.SetMapKey("a", a)
			_ = res.SetMapKey("b", b)
			return res
		},
		s: `merge(a, b)`,
		outFn: func() *model.Value {
			out := model.NewMapValue()
			_ = out.SetMapKey("key", model.NewStringValue("scalar"))
			return out
		},
	}.run)

	t.Run("deep three args", testCase{
		inFn: func() *model.Value {
			innerA := model.NewMapValue()
			_ = innerA.SetMapKey("x", model.NewIntValue(1))
			a := model.NewMapValue()
			_ = a.SetMapKey("nested", innerA)

			innerB := model.NewMapValue()
			_ = innerB.SetMapKey("y", model.NewIntValue(2))
			b := model.NewMapValue()
			_ = b.SetMapKey("nested", innerB)

			innerC := model.NewMapValue()
			_ = innerC.SetMapKey("z", model.NewIntValue(3))
			c := model.NewMapValue()
			_ = c.SetMapKey("nested", innerC)

			res := model.NewMapValue()
			_ = res.SetMapKey("a", a)
			_ = res.SetMapKey("b", b)
			_ = res.SetMapKey("c", c)
			return res
		},
		s: `merge(a, b, c)`,
		outFn: func() *model.Value {
			inner := model.NewMapValue()
			_ = inner.SetMapKey("x", model.NewIntValue(1))
			_ = inner.SetMapKey("y", model.NewIntValue(2))
			_ = inner.SetMapKey("z", model.NewIntValue(3))
			out := model.NewMapValue()
			_ = out.SetMapKey("nested", inner)
			return out
		},
	}.run)
}
