package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/dencoding"
	"github.com/tomwright/dasel/v3/model"
)

func TestUnary(t *testing.T) {
	t.Run("not", func(t *testing.T) {
		t.Run("literals", func(t *testing.T) {
			t.Run("not true", testCase{
				s:   `!true`,
				out: model.NewBoolValue(false),
			}.run)
			t.Run("not not true", testCase{
				s:   `!!true`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("not not not true", testCase{
				s:   `!!!true`,
				out: model.NewBoolValue(false),
			}.run)
			t.Run("not false", testCase{
				s:   `!false`,
				out: model.NewBoolValue(true),
			}.run)
			t.Run("not not false", testCase{
				s:   `!!false`,
				out: model.NewBoolValue(false),
			}.run)
			t.Run("not not not false", testCase{
				s:   `!!!false`,
				out: model.NewBoolValue(true),
			}.run)
		})
		t.Run("variables", func(t *testing.T) {
			in := func() *model.Value {
				return model.NewValue(dencoding.NewMap().
					Set("t", true).
					Set("f", false))
			}
			t.Run("not true", testCase{
				s:    `!t`,
				inFn: in,
				out:  model.NewBoolValue(false),
			}.run)
			t.Run("not not true", testCase{
				s:    `!!t`,
				inFn: in,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("not not not true", testCase{
				s:    `!!!t`,
				inFn: in,
				out:  model.NewBoolValue(false),
			}.run)
			t.Run("not false", testCase{
				s:    `!f`,
				inFn: in,
				out:  model.NewBoolValue(true),
			}.run)
			t.Run("not not false", testCase{
				s:    `!!f`,
				inFn: in,
				out:  model.NewBoolValue(false),
			}.run)
			t.Run("not not not false", testCase{
				s:    `!!!f`,
				inFn: in,
				out:  model.NewBoolValue(true),
			}.run)
		})
	})
}
