package execution_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/internal/ptr"
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
			if in == nil {
				in = model.NewValue(nil)
			}
			exp := tc.out
			if tc.outFn != nil {
				exp = tc.outFn()
			}
			res, err := execution.ExecuteSelector(tc.s, in)
			if err != nil {
				t.Fatal(err)
			}

			equal, err := res.EqualTypeValue(exp)
			if err != nil {
				t.Fatal(err)
			}
			if !equal {
				t.Errorf("unexpected type: %v", cmp.Diff(exp, res))
			}
		}
	}

	t.Run("binary expressions", func(t *testing.T) {
		t.Run("math", func(t *testing.T) {
			t.Run("literals", func(t *testing.T) {
				t.Run("addition", runTest(testCase{
					s:   `1 + 2`,
					out: model.NewIntValue(3),
				}))
				t.Run("subtraction", runTest(testCase{
					s:   `5 - 2`,
					out: model.NewIntValue(3),
				}))
				t.Run("multiplication", runTest(testCase{
					s:   `5 * 2`,
					out: model.NewIntValue(10),
				}))
				t.Run("division", runTest(testCase{
					s:   `10 / 2`,
					out: model.NewIntValue(5),
				}))
				t.Run("modulus", runTest(testCase{
					s:   `10 % 3`,
					out: model.NewIntValue(1),
				}))
				t.Run("ordering", runTest(testCase{
					s:   `45.2 + 5 * 4 - 2 / 2`, // 45.2 + (5 * 4) - (2 / 2) = 45.2 + 20 - 1 = 64.2
					out: model.NewFloatValue(64.2),
				}))
				t.Run("ordering with groups", runTest(testCase{
					s:   `(45.2 + 5) * ((4 - 2) / 2)`, // (45.2 + 5) * ((4 - 2) / 2) = (50.2) * ((2) / 2) = (50.2) * (1) = 50.2
					out: model.NewFloatValue(50.2),
				}))
			})
			t.Run("variables", func(t *testing.T) {
				in := func() *model.Value {
					return model.NewValue(dencoding.NewMap().
						Set("one", 1).
						Set("two", 2).
						Set("three", 3).
						Set("four", 4).
						Set("five", 5).
						Set("six", 6).
						Set("seven", 7).
						Set("eight", 8).
						Set("nine", 9).
						Set("ten", 10).
						Set("fortyfivepoint2", 45.2))
				}
				t.Run("addition", runTest(testCase{
					inFn: in,
					s:    `one + two`,
					out:  model.NewIntValue(3),
				}))
				t.Run("subtraction", runTest(testCase{
					inFn: in,
					s:    `five - two`,
					out:  model.NewIntValue(3),
				}))
				t.Run("multiplication", runTest(testCase{
					inFn: in,
					s:    `five * two`,
					out:  model.NewIntValue(10),
				}))
				t.Run("division", runTest(testCase{
					inFn: in,
					s:    `ten / two`,
					out:  model.NewIntValue(5),
				}))
				t.Run("modulus", runTest(testCase{
					inFn: in,
					s:    `ten % three`,
					out:  model.NewIntValue(1),
				}))
				t.Run("ordering", runTest(testCase{
					inFn: in,
					s:    `fortyfivepoint2 + five * four - two / two`, // 45.2 + (5 * 4) - (2 / 2) = 45.2 + 20 - 1 = 64.2
					out:  model.NewFloatValue(64.2),
				}))
				t.Run("ordering with groups", runTest(testCase{
					inFn: in,
					s:    `(fortyfivepoint2 + five) * ((four - two) / two)`, // (45.2 + 5) * ((4 - 2) / 2) = (50.2) * ((2) / 2) = (50.2) * (1) = 50.2
					out:  model.NewFloatValue(50.2),
				}))
			})
		})
		t.Run("comparison", func(t *testing.T) {
			t.Run("literals", func(t *testing.T) {
				t.Run("equal", runTest(testCase{
					s:   `1 == 1`,
					out: model.NewBoolValue(true),
				}))
				t.Run("not equal", runTest(testCase{
					s:   `1 != 1`,
					out: model.NewBoolValue(false),
				}))
				t.Run("greater than", runTest(testCase{
					s:   `2 > 1`,
					out: model.NewBoolValue(true),
				}))
				t.Run("greater than or equal", runTest(testCase{
					s:   `2 >= 2`,
					out: model.NewBoolValue(true),
				}))
				t.Run("less than", runTest(testCase{
					s:   `1 < 2`,
					out: model.NewBoolValue(true),
				}))
				t.Run("less than or equal", runTest(testCase{
					s:   `2 <= 2`,
					out: model.NewBoolValue(true),
				}))
			})

			t.Run("variables", func(t *testing.T) {
				in := func() *model.Value {
					return model.NewValue(dencoding.NewMap().
						Set("one", 1).
						Set("two", 2).
						Set("nested", dencoding.NewMap().
							Set("three", 3).
							Set("four", 4)))
				}
				t.Run("equal", runTest(testCase{
					inFn: in,
					s:    `one == one`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("not equal", runTest(testCase{
					inFn: in,
					s:    `one != one`,
					out:  model.NewBoolValue(false),
				}))
				t.Run("greater than", runTest(testCase{
					inFn: in,
					s:    `two > one`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("greater than or equal", runTest(testCase{
					inFn: in,
					s:    `two >= two`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("less than", runTest(testCase{
					inFn: in,
					s:    `one < two`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("less than or equal", runTest(testCase{
					inFn: in,
					s:    `two <= two`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("nested with math more than", runTest(testCase{
					inFn: in,
					s:    `nested.three + nested.four * 0 > one * 1`,
					out:  model.NewBoolValue(true),
				}))
				t.Run("nested with grouped math more than", runTest(testCase{
					inFn: in,
					s:    `(nested.three + nested.four) * 0 > one * 1`,
					out:  model.NewBoolValue(false),
				}))
			})
		})
	})

	t.Run("literal", func(t *testing.T) {
		t.Run("string", runTest(testCase{
			s:   `"hello"`,
			out: model.NewStringValue("hello"),
		}))
		t.Run("int", runTest(testCase{
			s:   `123`,
			out: model.NewIntValue(123),
		}))
		t.Run("float", runTest(testCase{
			s:   `123.4`,
			out: model.NewFloatValue(123.4),
		}))
		t.Run("true", runTest(testCase{
			s:   `true`,
			out: model.NewBoolValue(true),
		}))
		t.Run("false", runTest(testCase{
			s:   `false`,
			out: model.NewBoolValue(false),
		}))
	})

	t.Run("function", func(t *testing.T) {
		t.Run("add", func(t *testing.T) {
			t.Run("int", runTest(testCase{
				s:   `add(1, 2, 3)`,
				out: model.NewIntValue(6),
			}))
			t.Run("float", runTest(testCase{
				s:   `add(1f, 2.5, 3.5)`,
				out: model.NewFloatValue(7),
			}))
			t.Run("mixed", runTest(testCase{
				s:   `add(1, 2f)`,
				out: model.NewFloatValue(3),
			}))
			t.Run("properties", func(t *testing.T) {
				in := func() *model.Value {
					return model.NewValue(dencoding.NewMap().
						Set("numbers", dencoding.NewMap().
							Set("one", 1).
							Set("two", 2).
							Set("three", 3)).
						Set("nums", []any{1, 2, 3}))
				}
				t.Run("nested props", runTest(testCase{
					inFn: in,
					s:    `numbers.one + add(numbers.two, numbers.three)`,
					out:  model.NewIntValue(6),
				}))
				t.Run("add on end of chain", runTest(testCase{
					inFn: in,
					s:    `numbers.one + numbers.add(two, three)`,
					out:  model.NewIntValue(6),
				}))
				t.Run("add with map and spread on slice with $this addition and grouping", runTest(testCase{
					inFn: in,
					s:    `add(nums.map(($this + 1))...)`,
					out:  model.NewIntValue(9),
				}))
				t.Run("add with map and spread on slice with $this addition", runTest(testCase{
					inFn: in,
					s:    `add(nums.map($this + 1 - 2)...)`,
					out:  model.NewIntValue(3),
				}))
			})
		})
	})

	t.Run("get", func(t *testing.T) {
		inputMap := func() *model.Value {
			return model.NewValue(
				dencoding.NewMap().
					Set("title", "Mr").
					Set("age", int64(31)).
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
		t.Run("concat with grouping", runTest(testCase{
			in:  inputMap(),
			s:   `title + " " + (name.first) + " " + (name.last)`,
			out: model.NewStringValue("Mr Tom Wright"),
		}))
		t.Run("concat", runTest(testCase{
			in:  inputMap(),
			s:   `title + " " + name.first + " " + name.last`,
			out: model.NewStringValue("Mr Tom Wright"),
		}))
		t.Run("add evaluated fields", runTest(testCase{
			in: inputMap(),
			s:  `{..., over30 = age > 30}`,
			outFn: func() *model.Value {
				return model.NewValue(
					dencoding.NewMap().
						Set("title", "Mr").
						Set("age", int64(31)).
						Set("name", dencoding.NewMap().
							Set("first", "Tom").
							Set("last", "Wright")).
						Set("over30", true),
				)
			},
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
				map (
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
