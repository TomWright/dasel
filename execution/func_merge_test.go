package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestFuncMerge(t *testing.T) {
	t.Run("shallow", testCase{
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
}
