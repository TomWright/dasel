package dasel

import (
	"errors"
	"reflect"
	"testing"
)

func sameSlice(x, y []interface{}) bool {
	if len(x) != len(y) {
		return false
	}

	if reflect.DeepEqual(x, y) {
		return true
	}

	// Test for equality ignoring ordering
	diff := make([]interface{}, len(y))
	for k, v := range y {
		diff[k] = v
	}
	for _, xv := range x {
		for di, dv := range diff {
			if reflect.DeepEqual(xv, dv) {
				diff = append(diff[0:di], diff[di+1:]...)
				break
			}
		}
	}

	return len(diff) == 0
}

func selectTest(selector string, original interface{}, exp []interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		c := newSelectContext(original, selector)

		values, err := c.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got := values.Interfaces()
		if !sameSlice(exp, got) {
			t.Errorf("expected %v, got %v", exp, got)
			return
		}
	}
}

func selectTestErr(selector string, original interface{}, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		c := newSelectContext(original, selector)

		_, err := c.Run()

		if !errors.Is(err, expErr) {
			t.Errorf("expected error: %v, got %v", expErr, err)
			return
		}
	}
}

func TestContext_Step(t *testing.T) {
	step1 := &Step{index: 0}
	step2 := &Step{index: 1}
	c := &Context{
		steps: []*Step{
			step1, step2,
		},
	}
	expSteps := map[int]*Step{
		-1: nil,
		0:  step1,
		1:  step2,
		2:  nil,
	}

	for index, exp := range expSteps {
		got := c.Step(index)
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}

func TestContext_WithMetadata(t *testing.T) {
	c := (&Context{}).
		WithMetadata("x", 1).
		WithMetadata("y", 2)

	expMetadata := map[string]interface{}{
		"x": 1,
		"y": 2,
		"z": nil,
	}

	for index, exp := range expMetadata {
		got := c.Metadata(index)
		if exp != got {
			t.Errorf("expected %v, got %v", exp, got)
		}
	}
}
