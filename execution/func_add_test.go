package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/model"
)

func TestFuncAdd(t *testing.T) {
	t.Run("int", testCase{
		s:   `add(1, 2, 3)`,
		out: model.NewIntValue(6),
	}.run)
	t.Run("float", testCase{
		s:   `add(1f, 2.5, 3.5)`,
		out: model.NewFloatValue(7),
	}.run)
	t.Run("mixed", testCase{
		s:   `add(1, 2f)`,
		out: model.NewFloatValue(3),
	}.run)
	t.Run("properties", func(t *testing.T) {
		in := func() *model.Value {
			return model.NewValue(dencoding.NewMap().
				Set("numbers", dencoding.NewMap().
					Set("one", 1).
					Set("two", 2).
					Set("three", 3)).
				Set("nums", []any{1, 2, 3}))
		}
		t.Run("nested props", testCase{
			inFn: in,
			s:    `numbers.one + add(numbers.two, numbers.three)`,
			out:  model.NewIntValue(6),
		}.run)
		t.Run("add on end of chain", testCase{
			inFn: in,
			s:    `numbers.one + numbers.add(two, three)`,
			out:  model.NewIntValue(6),
		}.run)
		t.Run("add with map and spread on slice with $this addition and grouping", testCase{
			inFn: in,
			s:    `add(nums.map(($this + 1))...)`,
			out:  model.NewIntValue(9),
		}.run)
		t.Run("add with map and spread on slice with $this addition", testCase{
			inFn: in,
			s:    `add(nums.map($this + 1 - 2)...)`,
			out:  model.NewIntValue(3),
		}.run)
	})
}
