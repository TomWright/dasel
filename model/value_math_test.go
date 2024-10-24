package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestValue_Add(t *testing.T) {
	run := func(a, b *model.Value, exp *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			got, err := a.Add(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			eq, err := got.EqualTypeValue(exp)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !eq {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("int", run(model.NewIntValue(1), model.NewIntValue(2), model.NewIntValue(3)))
		t.Run("float", run(model.NewIntValue(1), model.NewFloatValue(2), model.NewFloatValue(3)))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("int", run(model.NewFloatValue(1), model.NewIntValue(2), model.NewFloatValue(3)))
		t.Run("float", run(model.NewFloatValue(1), model.NewFloatValue(2), model.NewFloatValue(3)))
	})
	t.Run("string", func(t *testing.T) {
		t.Run("string", run(model.NewStringValue("hello"), model.NewStringValue(" world"), model.NewStringValue("hello world")))
	})
}

func TestValue_Subtract(t *testing.T) {
	run := func(a, b *model.Value, exp *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			got, err := a.Subtract(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			eq, err := got.EqualTypeValue(exp)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !eq {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("int", run(model.NewIntValue(3), model.NewIntValue(2), model.NewIntValue(1)))
		t.Run("float", run(model.NewIntValue(3), model.NewFloatValue(2), model.NewFloatValue(1)))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("int", run(model.NewFloatValue(3), model.NewIntValue(2), model.NewFloatValue(1)))
		t.Run("float", run(model.NewFloatValue(3), model.NewFloatValue(2), model.NewFloatValue(1)))
	})
}

func TestValue_Multiply(t *testing.T) {
	run := func(a, b *model.Value, exp *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			got, err := a.Multiply(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			eq, err := got.EqualTypeValue(exp)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !eq {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("int", run(model.NewIntValue(3), model.NewIntValue(2), model.NewIntValue(6)))
		t.Run("float", run(model.NewIntValue(3), model.NewFloatValue(2), model.NewFloatValue(6)))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("int", run(model.NewFloatValue(3), model.NewIntValue(2), model.NewFloatValue(6)))
		t.Run("float", run(model.NewFloatValue(3), model.NewFloatValue(2), model.NewFloatValue(6)))
	})
}

func TestValue_Divide(t *testing.T) {
	run := func(a, b *model.Value, exp *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			got, err := a.Divide(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			eq, err := got.EqualTypeValue(exp)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !eq {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("int", run(model.NewIntValue(6), model.NewIntValue(2), model.NewIntValue(3)))
		t.Run("float", run(model.NewIntValue(6), model.NewFloatValue(2), model.NewFloatValue(3)))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("int", run(model.NewFloatValue(6), model.NewIntValue(2), model.NewFloatValue(3)))
		t.Run("float", run(model.NewFloatValue(6), model.NewFloatValue(2), model.NewFloatValue(3)))
	})
}

func TestValue_Modulo(t *testing.T) {
	run := func(a, b *model.Value, exp *model.Value) func(*testing.T) {
		return func(t *testing.T) {
			got, err := a.Modulo(b)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			eq, err := got.EqualTypeValue(exp)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}
			if !eq {
				t.Errorf("expected %v, got %v", exp, got)
			}
		}
	}
	t.Run("int", func(t *testing.T) {
		t.Run("int", run(model.NewIntValue(10), model.NewIntValue(3), model.NewIntValue(1)))
		t.Run("float", run(model.NewIntValue(10), model.NewFloatValue(3), model.NewFloatValue(1)))
	})
	t.Run("float", func(t *testing.T) {
		t.Run("int", run(model.NewFloatValue(10), model.NewIntValue(3), model.NewFloatValue(1)))
		t.Run("float", run(model.NewFloatValue(10), model.NewFloatValue(3), model.NewFloatValue(1)))
	})
}
