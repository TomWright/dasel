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
