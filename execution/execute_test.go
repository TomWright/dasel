package execution_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/ptr"
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

	t.Run("binary expressions", func(t *testing.T) {
		t.Run("math", func(t *testing.T) {
			t.Run("literals", func(t *testing.T) {
				t.Run("addition", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `1 + 2`,
					out: model.NewIntValue(3),
				}))
				t.Run("subtraction", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `5 - 2`,
					out: model.NewIntValue(3),
				}))
				t.Run("multiplication", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `5 * 2`,
					out: model.NewIntValue(10),
				}))
				t.Run("division", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `10 / 2`,
					out: model.NewIntValue(5),
				}))
				t.Run("modulus", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `10 % 3`,
					out: model.NewIntValue(1),
				}))
				t.Run("ordering", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `45.2 + 5 * 4 - 2 / 2`, // 45.2 + (5 * 4) - (2 / 2) = 45.2 + 20 - 1 = 64.2
					out: model.NewFloatValue(64.2),
				}))
				t.Run("ordering with groups", runTest(testCase{
					in:  model.NewValue(nil),
					s:   `(45.2 + 5) * ((4 - 2) / 2)`, // (45.2 + 5) * ((4 - 2) / 2) = (50.2) * ((2) / 2) = (50.2) * (1) = 50.2
					out: model.NewFloatValue(50.2),
				}))
			})
		})
		t.Run("comparison", func(t *testing.T) {
			t.Run("equal", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `1 == 1`,
				out: model.NewBoolValue(true),
			}))
			t.Run("not equal", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `1 != 1`,
				out: model.NewBoolValue(false),
			}))
			t.Run("greater than", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `2 > 1`,
				out: model.NewBoolValue(true),
			}))
			t.Run("greater than or equal", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `2 >= 2`,
				out: model.NewBoolValue(true),
			}))
			t.Run("less than", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `1 < 2`,
				out: model.NewBoolValue(true),
			}))
			t.Run("less than or equal", runTest(testCase{
				in:  model.NewValue(nil),
				s:   `2 <= 2`,
				out: model.NewBoolValue(true),
			}))
		})
	})

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
				return model.NewValue([]any{
					dencoding.NewMap().Set("number", 1),
					dencoding.NewMap().Set("number", 2),
					dencoding.NewMap().Set("number", 3),
				})
			},
			s: `map(number)`,
			outFn: func() *model.Value {
				return model.NewValue([]any{1, 2, 3})
			},
		}))
		t.Run("with chain of selectors", runTest(testCase{
			inFn: func() *model.Value {
				return model.NewValue([]any{
					dencoding.NewMap().Set("foo", 1).Set("bar", 4),
					dencoding.NewMap().Set("foo", 2).Set("bar", 5),
					dencoding.NewMap().Set("foo", 3).Set("bar", 6),
				})
			},
			s: `
				.map (
					{
						total = add( foo, bar, 1 )
					}
				)
				.map ( total )`,
			outFn: func() *model.Value {
				return model.NewValue([]any{ptr.To(int64(6)), ptr.To(int64(8)), ptr.To(int64(10))})
			},
		}))
	})
}
