package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/model"
)

func TestObject(t *testing.T) {
	inputMap := func() *model.Value {
		return model.NewValue(dencoding.NewMap().
			Set("title", "Mr").
			Set("age", int64(30)).
			Set("name", dencoding.NewMap().
				Set("first", "Tom").
				Set("last", "Wright")))
	}
	t.Run("get", testCase{
		in: inputMap(),
		s:  `{title}`,
		outFn: func() *model.Value {
			return model.NewValue(dencoding.NewMap().Set("title", "Mr"))
			//res := model.NewMapValue()
			//_ = res.SetMapKey("title", model.NewStringValue("Mr"))
			//return res
		},
	}.run)
	t.Run("get multiple", testCase{
		in: inputMap(),
		s:  `{title, age}`,
		outFn: func() *model.Value {
			return model.NewValue(dencoding.NewMap().Set("title", "Mr").Set("age", int64(30)))
			//res := model.NewMapValue()
			//_ = res.SetMapKey("title", model.NewStringValue("Mr"))
			//_ = res.SetMapKey("age", model.NewIntValue(30))
			//return res
		},
	}.run)
	t.Run("get with spread", testCase{
		in: inputMap(),
		s:  `{...}`,
		outFn: func() *model.Value {
			res := inputMap()
			return res
		},
	}.run)
	t.Run("set", testCase{
		in: inputMap(),
		s:  `{title:"Mrs"}`,
		outFn: func() *model.Value {
			res := model.NewMapValue()
			_ = res.SetMapKey("title", model.NewStringValue("Mrs"))
			return res
		},
	}.run)
	t.Run("set with spread", testCase{
		in: inputMap(),
		s:  `{..., title:"Mrs"}`,
		outFn: func() *model.Value {
			res := inputMap()
			_ = res.SetMapKey("title", model.NewStringValue("Mrs"))
			return res
		},
	}.run)
	t.Run("merge with spread", testCase{
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
		s: `{a..., b..., x: 1}`,
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
			if err := b.SetMapKey("x", model.NewIntValue(1)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			return b
		},
	}.run)
}
