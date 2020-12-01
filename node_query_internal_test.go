package dasel

import (
	"reflect"
	"testing"
)

func assertErrResult(t *testing.T, expErr error, gotErr error) bool {
	if expErr == nil && gotErr != nil {
		t.Errorf("expected err %v, got %v", expErr, gotErr)
		return false
	}
	if expErr != nil && gotErr == nil {
		t.Errorf("expected err %v, got %v", expErr, gotErr)
		return false
	}
	if expErr != nil && gotErr != nil && gotErr.Error() != expErr.Error() {
		t.Errorf("expected err %v, got %v", expErr, gotErr)
		return false
	}
	return true
}

func assertQueryResult(t *testing.T, exp reflect.Value, expErr error, got reflect.Value, gotErr error) bool {
	if !assertErrResult(t, expErr, gotErr) {
		return false
	}
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("expected result %v, got %v", exp, got)
		return false
	}
	return true
}

func getNodeWithValue(value interface{}) *Node {
	rootNode := New(value)
	nextNode := &Node{
		Previous: rootNode,
	}
	rootNode.Next = nextNode
	return nextNode
}

func TestFindValueProperty(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "."
		got, err := findValueProperty(n, false)
		assertQueryResult(t, nilValue(), &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		n := getNodeWithValue(map[string]interface{}{})
		n.Selector.Current = "x"
		got, err := findValueProperty(n, false)
		assertQueryResult(t, nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}, got, err)
	})
}

func TestFindValueIndex(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "."
		got, err := findValueIndex(n, false)
		assertQueryResult(t, nilValue(), &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Current = "[0]"
		n.Selector.Index = 0
		got, err := findValueIndex(n, false)
		assertQueryResult(t, nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		val := map[string]interface{}{}
		n := getNodeWithValue(val)
		n.Selector.Current = "[0]"
		n.Selector.Index = 0
		got, err := findValueIndex(n, false)
		assertQueryResult(t, nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: reflect.TypeOf(val).Kind()}, got, err)
	})
}

func TestFindValueNextAvailableIndex(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Current = "[0]"
		n.Selector.Index = 0
		got, err := findNextAvailableIndex(n, false)
		assertQueryResult(t, nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}, got, err)
	})
}

func TestFindValueDynamic(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		n := getNodeWithValue(nil)
		n.Previous.Selector.Current = "."
		got, err := findValueDynamic(n, false)
		assertQueryResult(t, nilValue(), &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Current = "(name=x)"
		n.Selector.Conditions = []Condition{
			&EqualCondition{Key: "name", Value: "x"},
		}
		got, err := findValueDynamic(n, false)
		assertQueryResult(t, nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}, got, err)
	})
	t.Run("NotFoundWithCreate", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Current = "(name=x)"
		n.Selector.Conditions = []Condition{
			&EqualCondition{Key: "name", Value: "x"},
		}
		got, err := findValueDynamic(n, true)
		if !assertQueryResult(t, nilValue(), nil, got, err) {
			return
		}
		if exp, got := "NEXT_AVAILABLE_INDEX", n.Selector.Type; exp != got {
			t.Errorf("expected type of %s, got %s", exp, got)
			return
		}
	})
	t.Run("UnsupportedCheckType", func(t *testing.T) {
		itemVal := 1
		val := []interface{}{
			itemVal,
		}
		n := getNodeWithValue(val)
		n.Selector.Current = "(name=x)"
		n.Selector.Conditions = []Condition{
			&EqualCondition{Key: "name", Value: "x"},
		}
		got, err := findValueDynamic(n, false)
		assertQueryResult(t, nilValue(), &UnhandledCheckType{Value: reflect.TypeOf(itemVal).Kind().String()}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		val := 0
		n := getNodeWithValue(val)
		n.Selector.Current = "(name=x)"
		n.Selector.Conditions = []Condition{
			&EqualCondition{Key: "name", Value: "x"},
		}
		got, err := findValueDynamic(n, false)
		assertQueryResult(t, nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: reflect.TypeOf(val).Kind()}, got, err)
	})
}

func TestFindValue(t *testing.T) {
	t.Run("MissingPreviousNode", func(t *testing.T) {
		n := New(nil)
		got, err := findValue(n, false)
		assertQueryResult(t, nilValue(), ErrMissingPreviousNode, got, err)
	})
	t.Run("UnsupportedSelector", func(t *testing.T) {
		n := getNodeWithValue([]interface{}{})
		n.Selector.Raw = "BAD"
		got, err := findValue(n, false)
		assertQueryResult(t, nilValue(), &UnsupportedSelector{Selector: "BAD"}, got, err)
	})
}
