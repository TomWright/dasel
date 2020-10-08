package dasel

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPropagateValueProperty(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "."
		err := propagateValueProperty(n)
		assertErrResult(t, &UnexpectedPreviousNilValue{Selector: "."}, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		val := make([]interface{}, 0)
		n := getNodeWithValue(val)
		err := propagateValueProperty(n)
		assertErrResult(t, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}, err)
	})
}

func TestPropagateValueIndex(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "."
		err := propagateValueIndex(n)
		assertErrResult(t, &UnexpectedPreviousNilValue{Selector: "."}, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		val := map[string]interface{}{}
		n := getNodeWithValue(val)
		n.Selector.Current = "[0]"
		n.Selector.Index = 0
		err := propagateValueIndex(n)
		assertErrResult(t, &UnsupportedTypeForSelector{Selector: n.Selector, Value: reflect.TypeOf(val).Kind()}, err)
	})
	t.Run("ExistingIndex", func(t *testing.T) {
		val := []interface{}{
			"hello",
		}
		n := getNodeWithValue(val)
		n.Selector.Current = "[0]"
		n.Selector.Index = 0
		n.Value = reflect.ValueOf("world")
		err := propagateValueIndex(n)
		if !assertErrResult(t, nil, err) {
			return
		}
		if !reflect.DeepEqual(val[0], "world") {
			t.Errorf("expected val %v, got %v", "world", val[0])
		}
	})
}

func TestPropagateValueNextAvailableIndex(t *testing.T) {
	t.Run("MissingPreviousNode", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "x"
		err := propagateValueNextAvailableIndex(n)
		assertErrResult(t, &UnexpectedPreviousNilValue{Selector: "x"}, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		val := map[string]interface{}{}
		n := getNodeWithValue(val)
		err := propagateValueNextAvailableIndex(n)
		assertErrResult(t, &UnsupportedTypeForSelector{Selector: n.Selector, Value: reflect.TypeOf(val).Kind()}, err)
	})
}

func TestPropagateValue(t *testing.T) {
	t.Run("MissingPreviousNode", func(t *testing.T) {
		n := New(nil)
		err := propagateValue(n)
		assertErrResult(t, nil, err)
	})
	t.Run("UnsupportedSelector", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Type = "BAD"
		err := propagateValue(n)
		assertErrResult(t, &UnsupportedSelector{Selector: "BAD"}, err)
	})
}

func TestPropagate(t *testing.T) {
	t.Run("UnsupportedSelector", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Type = "BAD"
		err := propagate(n)
		assertErrResult(t, fmt.Errorf("could not propagate value: %w", &UnsupportedSelector{Selector: "BAD"}), err)
	})
}
