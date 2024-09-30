package execution_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestExecuteSelector_HappyPath(t *testing.T) {
	type testCase struct {
		in    *model.Value
		inFn  func() *model.Value
		s     string
		out   *model.Value
		outFn func() *model.Value
	}

	runTest := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			in := tc.in
			if tc.inFn != nil {
				in = tc.inFn()
			}
			exp := tc.out
			if tc.outFn != nil {
				exp = tc.outFn()
			}
			res, err := execution.ExecuteSelector(tc.s, in)
			if err != nil {
				t.Fatal(err)
			}

			toInterface := func(v *model.Value) interface{} {
				if v == nil {
					return nil
				}
				if v.IsMap() {
					m, _ := v.MapValue()
					return m.KeyValues()
				}
				return v.UnpackKinds(reflect.Ptr).Interface()
			}
			expV, gotV := toInterface(exp), toInterface(res)

			if !cmp.Equal(expV, gotV, cmpopts.IgnoreUnexported(dencoding.Map{})) {
				t.Errorf("unexpected result: %v", cmp.Diff(expV, gotV))
			}
		}
	}

	t.Run("literal", func(t *testing.T) {
		t.Run("string", runTest(testCase{
			in:  model.NewValue(nil),
			s:   `"hello"`,
			out: model.NewStringValue("hello"),
		}))
		t.Run("int", runTest(testCase{
			in:  model.NewValue(nil),
			s:   `123`,
			out: model.NewIntValue(123),
		}))
		t.Run("float", runTest(testCase{
			in:  model.NewValue(nil),
			s:   `123.4`,
			out: model.NewFloatValue(123.4),
		}))
		t.Run("true", runTest(testCase{
			in:  model.NewValue(nil),
			s:   `true`,
			out: model.NewBoolValue(true),
		}))
		t.Run("false", runTest(testCase{
			in:  model.NewValue(nil),
			s:   `false`,
			out: model.NewBoolValue(false),
		}))
	})

	t.Run("function", func(t *testing.T) {
		t.Run("add", func(t *testing.T) {
			t.Run("int", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `add(1, 2, 3)`,
				out: model.NewIntValue(6),
			}))
			t.Run("float", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `add(1f, 2.5, 3.5)`,
				out: model.NewFloatValue(7),
			}))
			t.Run("mixed", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `add(1, 2f)`,
				out: model.NewFloatValue(3),
			}))
		})
	})

	t.Run("get", func(t *testing.T) {
		inputMap := func() *model.Value {
			return model.NewValue(
				dencoding.NewMap().
					Set("title", "Mr").
					Set("name", dencoding.NewMap().
						Set("first", "Tom").
						Set("last", "Wright")),
			)
		}
		t.Run("property", runTest(testCase{
			in:  inputMap(),
			s:   `title`,
			out: model.NewStringValue("Mr"),
		}))
		t.Run("nested property", runTest(testCase{
			in:  inputMap(),
			s:   `name.first`,
			out: model.NewStringValue("Tom"),
		}))
	})

	t.Run("object", func(t *testing.T) {
		inputMap := func() *model.Value {
			return model.NewValue(dencoding.NewMap().
				Set("title", "Mr").
				Set("age", int64(30)).
				Set("name", dencoding.NewMap().
					Set("first", "Tom").
					Set("last", "Wright")))
		}
		t.Run("get", runTest(testCase{
			in: inputMap(),
			s:  `{title}`,
			outFn: func() *model.Value {
				return model.NewValue(dencoding.NewMap().Set("title", "Mr"))
				//res := model.NewMapValue()
				//_ = res.SetMapKey("title", model.NewStringValue("Mr"))
				//return res
			},
		}))
		t.Run("get multiple", runTest(testCase{
			in: inputMap(),
			s:  `{title, age}`,
			outFn: func() *model.Value {
				return model.NewValue(dencoding.NewMap().Set("title", "Mr").Set("age", int64(30)))
				//res := model.NewMapValue()
				//_ = res.SetMapKey("title", model.NewStringValue("Mr"))
				//_ = res.SetMapKey("age", model.NewIntValue(30))
				//return res
			},
		}))
		t.Run("get with spread", runTest(testCase{
			in: inputMap(),
			s:  `{...}`,
			outFn: func() *model.Value {
				res := inputMap()
				return res
			},
		}))
		t.Run("set", runTest(testCase{
			in: inputMap(),
			s:  `{title="Mrs"}`,
			outFn: func() *model.Value {
				res := model.NewMapValue()
				_ = res.SetMapKey("title", model.NewStringValue("Mrs"))
				return res
			},
		}))
		t.Run("set with spread", runTest(testCase{
			in: inputMap(),
			s:  `{..., title="Mrs"}`,
			outFn: func() *model.Value {
				res := inputMap()
				_ = res.SetMapKey("title", model.NewStringValue("Mrs"))
				return res
			},
		}))
	})

	t.Run("map", func(t *testing.T) {
		t.Run("property from slice of maps", runTest(testCase{
			inFn: func() *model.Value {
				r := model.NewSliceValue()

				m1 := model.NewMapValue()
				_ = m1.SetMapKey("number", model.NewIntValue(1))
				m2 := model.NewMapValue()
				_ = m2.SetMapKey("number", model.NewIntValue(2))
				m3 := model.NewMapValue()
				_ = m3.SetMapKey("number", model.NewIntValue(3))

				_ = r.Append(m1)
				_ = r.Append(m2)
				_ = r.Append(m3)

				return r
			},
			s: `map(number)`,
			outFn: func() *model.Value {
				r := model.NewSliceValue()
				_ = r.Append(model.NewIntValue(1))
				_ = r.Append(model.NewIntValue(2))
				_ = r.Append(model.NewIntValue(3))
				return r
			},
		}))
	})
}
