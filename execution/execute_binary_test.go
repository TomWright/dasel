package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/model/orderedmap"
)

func TestBinary(t *testing.T) {
	t.Run("math", func(t *testing.T) {
		t.Run("literals", func(t *testing.T) {
			t.Run("addition", testCase{
				s:   `1 + 2`,
				out: model.NewIntValue(3),
			}.run)
			t.Run("subtraction", testCase{
				s:   `5 - 2`,
				out: model.NewIntValue(3),
			}.run)
			t.Run("multiplication", testCase{
				s:   `5 * 2`,
				out: model.NewIntValue(10),
			}.run)
			t.Run("division", testCase{
				s:   `10 / 2`,
				out: model.NewIntValue(5),
			}.run)
			t.Run("modulus", testCase{
				s:   `10 % 3`,
				out: model.NewIntValue(1),
			}.run)
			t.Run("ordering", testCase{
				s:   `45.2 + 5 * 4 - 2 / 2`, // 45.2 + (5 * 4) - (2 / 2) = 45.2 + 20 - 1 = 64.2
				out: model.NewFloatValue(64.2),
			}.run)
			t.Run("ordering with groups", testCase{
				s:   `(45.2 + 5) * ((4 - 2) / 2)`, // (45.2 + 5) * ((4 - 2) / 2) = (50.2) * ((2) / 2) = (50.2) * (1) = 50.2
				out: model.NewFloatValue(50.2),
			}.run)
			t.Run("ordering with groups", testCase{
				s:   `1 + 1 - 1 + 1 * 2`, // 1 + 1 - 1 + (1 * 2) = 1 + 1 - 1 + 2 = 3
				out: model.NewIntValue(3),
			}.run)
		})
		t.Run("variables", func(t *testing.T) {
			in := func() *model.Value {
				return model.NewValue(orderedmap.NewMap().
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
			t.Run("addition", testCase{
				inFn: in,
				s:    `one + two`,
				out:  model.NewIntValue(3),
			}.run)
			t.Run("subtraction", testCase{
				inFn: in,
				s:    `five - two`,
				out:  model.NewIntValue(3),
			}.run)
			t.Run("multiplication", testCase{
				inFn: in,
				s:    `five * two`,
				out:  model.NewIntValue(10),
			}.run)
			t.Run("division", testCase{
				inFn: in,
				s:    `ten / two`,
				out:  model.NewIntValue(5),
			}.run)
			t.Run("modulus", testCase{
				inFn: in,
				s:    `ten % three`,
				out:  model.NewIntValue(1),
			}.run)
			t.Run("ordering", testCase{
				inFn: in,
				s:    `fortyfivepoint2 + five * four - two / two`, // 45.2 + (5 * 4) - (2 / 2) = 45.2 + 20 - 1 = 64.2
				out:  model.NewFloatValue(64.2),
			}.run)
			t.Run("ordering with groups", testCase{
				inFn: in,
				s:    `(fortyfivepoint2 + five) * ((four - two) / two)`, // (45.2 + 5) * ((4 - 2) / 2) = (50.2) * ((2) / 2) = (50.2) * (1) = 50.2
				out:  model.NewFloatValue(50.2),
			}.run)
		})
	})
	t.Run("comparison", func(t *testing.T) {
		t.Run("literals", func(t *testing.T) {
			t.Run("equal", testCase{
				s:   `1 == 1`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("not equal", testCase{
				s:   `1 != 1`,
				out: model.NewBoolValue(false),
			}.run)
			t.Run("greater than", testCase{
				s:   `2 > 1`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("greater than or equal", testCase{
				s:   `2 >= 2`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("less than", testCase{
				s:   `1 < 2`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("less than or equal", testCase{
				s:   `2 <= 2`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("like", testCase{
				s:   `"hello world" =~ r/ello/`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("not like", testCase{
				s:   `"hello world" !~ r/helloworld/`,
				out: model.NewBoolValue(true),
			}.run)
		})

		t.Run("variables", func(t *testing.T) {
			in := func() *model.Value {
				return model.NewValue(orderedmap.NewMap().
					Set("one", 1).
					Set("two", 2).
					Set("nested", orderedmap.NewMap().
						Set("three", 3).
						Set("four", 4)))
			}
			t.Run("equal", testCase{
				inFn: in,
				s:    `one == one`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("not equal", testCase{
				inFn: in,
				s:    `one != one`,
				out:  model.NewBoolValue(false),
			}.run)
			t.Run("greater than", testCase{
				inFn: in,
				s:    `two > one`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("greater than or equal", testCase{
				inFn: in,
				s:    `two >= two`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("less than", testCase{
				inFn: in,
				s:    `one < two`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("less than or equal", testCase{
				inFn: in,
				s:    `two <= two`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("nested with math more than", testCase{
				inFn: in,
				s:    `nested.three + nested.four * 0 > one * 1`,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("nested with grouped math more than", testCase{
				inFn: in,
				s:    `(nested.three + nested.four) * 0 > one * 1`,
				out:  model.NewBoolValue(false),
			}.run)
		})

		t.Run("coalesce", func(t *testing.T) {
			t.Run("literals", func(t *testing.T) {
				t.Run("coalesce", testCase{
					s:   `null ?? 1`,
					out: model.NewIntValue(1),
				}.run)
				t.Run("coalesce with null", testCase{
					s:   `null ?? null`,
					out: model.NewNullValue(),
				}.run)
				t.Run("coalesce with null and value", testCase{
					s:   `null ?? 2`,
					out: model.NewIntValue(2),
				}.run)
				t.Run("coalesce with value", testCase{
					s:   `1 ?? 2`,
					out: model.NewIntValue(1),
				}.run)
			})
			t.Run("variables", func(t *testing.T) {
				in := func() *model.Value {
					return model.NewValue(orderedmap.NewMap().
						Set("one", 1).
						Set("two", 2).
						Set("nested", orderedmap.NewMap().
							Set("one", 1).
							Set("two", 2).
							Set("three", 3).
							Set("four", 4)).
						Set("list", []any{1, 2, 3}))
				}
				t.Run("coalesce", testCase{
					inFn: in,
					s:    `nested.five ?? one`,
					out:  model.NewIntValue(1),
				}.run)
				t.Run("coalesce with null", testCase{
					inFn: in,
					s:    `nested.five ?? null`,
					out:  model.NewNullValue(),
				}.run)
				t.Run("coalesce with null and value", testCase{
					inFn: in,
					s:    `nested.five ?? 2`,
					out:  model.NewIntValue(2),
				}.run)
				t.Run("coalesce with value", testCase{
					inFn: in,
					s:    `nested.three ?? 2`,
					out:  model.NewIntValue(3),
				}.run)
				t.Run("coalesce with bad map key", testCase{
					inFn: in,
					s:    `nope ?? 2`,
					out:  model.NewIntValue(2),
				}.run)
				t.Run("coalesce with nested bad map key", testCase{
					inFn: in,
					s:    `nested.nope ?? 2`,
					out:  model.NewIntValue(2),
				}.run)
				t.Run("coalesce with list index", testCase{
					inFn: in,
					s:    `list[1] ?? 5`,
					out:  model.NewIntValue(2),
				}.run)
				t.Run("coalesce with list bad index", testCase{
					inFn: in,
					s:    `list[3] ?? 5`,
					out:  model.NewIntValue(5),
				}.run)
				t.Run("chained coalesce execute left to right", func(t *testing.T) {
					// These tests ensure the coalesces run in order.
					t.Run("no match", testCase{
						inFn: in,
						s:    `nested.five ?? nested.six ?? nested.seven ?? 10`,
						out:  model.NewIntValue(10),
					}.run)
					t.Run("first match when all exist", testCase{
						inFn: in,
						s:    `nested.one ?? nested.two ?? nested.three ?? 10`,
						out:  model.NewIntValue(1),
					}.run)
					t.Run("second match", testCase{
						inFn: in,
						s:    `nested.five ?? nested.two ?? nested.three ?? 10`,
						out:  model.NewIntValue(2),
					}.run)
					t.Run("third match", testCase{
						inFn: in,
						s:    `nested.five ?? nested.six ?? nested.three ?? 10`,
						out:  model.NewIntValue(3),
					}.run)
				})
			})
		})
	})
}
