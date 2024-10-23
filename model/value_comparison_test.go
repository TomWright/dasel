package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

type compareTestCase struct {
	a   *model.Value
	b   *model.Value
	exp bool
}

func TestValue_Equal(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.Equal(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("string", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("hello"),
			b:   model.NewStringValue("hello"),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewStringValue("hello"),
			b:   model.NewStringValue("world"),
			exp: false,
		}))
	})
	t.Run("int", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: false,
		}))
		t.Run("equal float", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: true,
		}))
		t.Run("not equal float", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: false,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: false,
		}))
		t.Run("equal int", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("not equal int", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
	})
	t.Run("bool", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewBoolValue(true),
			b:   model.NewBoolValue(true),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewBoolValue(true),
			b:   model.NewBoolValue(false),
			exp: false,
		}))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			b: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			b: model.NewValue(map[string]interface{}{
				"hello": "world2",
			}),
			exp: false,
		}))
	})
	t.Run("array", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a: model.NewValue([]interface{}{
				"hello", "world",
			}),
			b: model.NewValue([]interface{}{
				"hello", "world",
			}),
			exp: true,
		}))
		t.Run("not equal", run(compareTestCase{
			a: model.NewValue([]interface{}{
				"hello", "world",
			}),
			b: model.NewValue([]interface{}{
				"hello", "world2",
			}),
			exp: false,
		}))
	})
	t.Run("null", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewValue(nil),
			b:   model.NewValue(nil),
			exp: true,
		}))
	})
}

func TestValue_NotEqual(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.NotEqual(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("string", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("hello"),
			b:   model.NewStringValue("hello"),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewStringValue("hello"),
			b:   model.NewStringValue("world"),
			exp: true,
		}))
	})
	t.Run("int", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: true,
		}))
		t.Run("equal float", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: false,
		}))
		t.Run("not equal float", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: true,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: true,
		}))
		t.Run("equal int", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("not equal int", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
	})
	t.Run("bool", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewBoolValue(true),
			b:   model.NewBoolValue(true),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a:   model.NewBoolValue(true),
			b:   model.NewBoolValue(false),
			exp: true,
		}))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			b: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a: model.NewValue(map[string]interface{}{
				"hello": "world",
			}),
			b: model.NewValue(map[string]interface{}{
				"hello": "world2",
			}),
			exp: true,
		}))
	})
	t.Run("array", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a: model.NewValue([]interface{}{
				"hello", "world",
			}),
			b: model.NewValue([]interface{}{
				"hello", "world",
			}),
			exp: false,
		}))
		t.Run("not equal", run(compareTestCase{
			a: model.NewValue([]interface{}{
				"hello", "world",
			}),
			b: model.NewValue([]interface{}{
				"hello", "world2",
			}),
			exp: true,
		}))
	})
	t.Run("null", func(t *testing.T) {
		t.Run("equal", run(compareTestCase{
			a:   model.NewValue(nil),
			b:   model.NewValue(nil),
			exp: false,
		}))
	})
}

func TestValue_LessThan(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.LessThan(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(1.2),
			b:   model.NewFloatValue(1.1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: false,
		}))
	})
	t.Run("int float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewFloatValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: false,
		}))
	})
	t.Run("float int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(2),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("b"),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewStringValue("b"),
			b:   model.NewStringValue("a"),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("a"),
			exp: false,
		}))
	})
}

func TestValue_LessThanOrEqual(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.LessThanOrEqual(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(1.2),
			b:   model.NewFloatValue(1.1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: true,
		}))
	})
	t.Run("int float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewFloatValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: true,
		}))
	})
	t.Run("float int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(2),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(2),
			b:   model.NewIntValue(1),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("b"),
			exp: true,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewStringValue("b"),
			b:   model.NewStringValue("a"),
			exp: false,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("a"),
			exp: true,
		}))
	})
}

func TestValue_GreaterThan(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.GreaterThan(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(1.2),
			b:   model.NewFloatValue(1.1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: false,
		}))
	})
	t.Run("int float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewFloatValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: false,
		}))
	})
	t.Run("float int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(2),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: false,
		}))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("b"),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewStringValue("b"),
			b:   model.NewStringValue("a"),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("a"),
			exp: false,
		}))
	})
}

func TestValue_GreaterThanOrEqual(t *testing.T) {
	run := func(tc compareTestCase) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := tc.a.GreaterThanOrEqual(tc.b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			gotBool, err := got.BoolValue()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if gotBool != tc.exp {
				t.Errorf("expected %v, got %v", tc.exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(1.2),
			b:   model.NewFloatValue(1.1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewFloatValue(1.1),
			exp: true,
		}))
	})
	t.Run("int float", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewIntValue(2),
			b:   model.NewFloatValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewIntValue(1),
			b:   model.NewFloatValue(1),
			exp: true,
		}))
	})
	t.Run("float int", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewFloatValue(1.1),
			b:   model.NewIntValue(2),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewFloatValue(2),
			b:   model.NewIntValue(1),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewFloatValue(1),
			b:   model.NewIntValue(1),
			exp: true,
		}))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("less", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("b"),
			exp: false,
		}))
		t.Run("greater", run(compareTestCase{
			a:   model.NewStringValue("b"),
			b:   model.NewStringValue("a"),
			exp: true,
		}))
		t.Run("equal", run(compareTestCase{
			a:   model.NewStringValue("a"),
			b:   model.NewStringValue("a"),
			exp: true,
		}))
	})
}

func TestValue_Compare(t *testing.T) {
	run := func(a *model.Value, b *model.Value, exp int) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := a.Compare(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if got != exp {
				t.Errorf("expected %d, got %d", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("less", run(
			model.NewIntValue(1),
			model.NewIntValue(2),
			-1,
		))
		t.Run("greater", run(
			model.NewIntValue(2),
			model.NewIntValue(1),
			1,
		))
		t.Run("equal", run(
			model.NewIntValue(1),
			model.NewIntValue(1),
			0,
		))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("less", run(
			model.NewFloatValue(1.1),
			model.NewFloatValue(1.2),
			-1,
		))
		t.Run("greater", run(
			model.NewFloatValue(1.2),
			model.NewFloatValue(1.1),
			1,
		))
		t.Run("equal", run(
			model.NewFloatValue(1.1),
			model.NewFloatValue(1.1),
			0,
		))
	})
	t.Run("int float", func(t *testing.T) {
		t.Run("less", run(
			model.NewIntValue(1),
			model.NewFloatValue(2),
			-1,
		))
		t.Run("greater", run(
			model.NewIntValue(2),
			model.NewFloatValue(1),
			1,
		))
		t.Run("equal", run(
			model.NewIntValue(1),
			model.NewFloatValue(1),
			0,
		))
	})
	t.Run("float int", func(t *testing.T) {
		t.Run("less", run(
			model.NewFloatValue(1.1),
			model.NewIntValue(2),
			-1,
		))
		t.Run("greater", run(
			model.NewFloatValue(1.1),
			model.NewIntValue(1),
			1,
		))
		t.Run("equal", run(
			model.NewFloatValue(1),
			model.NewIntValue(1),
			0,
		))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("less", run(
			model.NewStringValue("a"),
			model.NewStringValue("b"),
			-1,
		))
		t.Run("greater", run(
			model.NewStringValue("b"),
			model.NewStringValue("a"),
			1,
		))
		t.Run("equal", run(
			model.NewStringValue("a"),
			model.NewStringValue("a"),
			0,
		))
	})
}
