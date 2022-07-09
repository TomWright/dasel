package dasel

import (
	"fmt"
	"reflect"
	"testing"
)

func assertQueryMultipleResult(t *testing.T, exp []reflect.Value, expErr error, got []*Node, gotErr error) bool {
	if !assertErrResult(t, expErr, gotErr) {
		return false
	}
	gotVals := make([]interface{}, len(got))
	if len(got) > 0 {
		for i, n := range got {
			if !n.Value.IsValid() {
				gotVals[i] = nil
				continue
			}
			gotVals[i] = n.Value.Interface()
		}
	}
	expVals := make([]interface{}, len(exp))
	if len(exp) > 0 {
		for i, n := range exp {
			if !n.IsValid() {
				expVals[i] = nil
				continue
			}
			expVals[i] = n.Interface()
		}
	}
	if !reflect.DeepEqual(expVals, gotVals) {
		t.Errorf("expected result %v, got %v", expVals, gotVals)
		return false
	}
	return true
}

func assertQueryMultipleResultOneOf(t *testing.T, exp [][]reflect.Value, expErr error, got []*Node, gotErr error) bool {
	if !assertErrResult(t, expErr, gotErr) {
		return false
	}
	gotVals := make([]interface{}, len(got))
	if len(got) > 0 {
		for i, n := range got {
			if !n.Value.IsValid() {
				gotVals[i] = nil
				continue
			}
			gotVals[i] = n.Value.Interface()
		}
	}

	for _, exp := range exp {
		expVals := make([]interface{}, len(exp))
		if len(exp) > 0 {
			for i, n := range exp {
				if !n.IsValid() {
					expVals[i] = nil
					continue
				}
				expVals[i] = n.Interface()
			}
		}

		if reflect.DeepEqual(expVals, gotVals) {
			return true
		}
	}

	t.Errorf("unexpected result: %v", gotVals)
	return false
}

func TestFindNodesProperty(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: ".", Raw: "."}
		got, err := findNodesProperty(selector, nilValue(), false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		previousValue := reflect.ValueOf(map[string]interface{}{})
		selector := Selector{Current: "x"}
		got, err := findNodesProperty(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		previousValue := reflect.ValueOf(0)
		selector := Selector{Current: "x"}
		got, err := findNodesProperty(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue}, got, err)
	})
}

func TestFindNodesLength(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		got, err := findNodesLength(selector, nilValue())
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: ".[#]"}, got, err)
	})
	t.Run("UnsupportedTypeInt", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := 0
		got, err := findNodesLength(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: reflect.ValueOf(val)}, got, err)
	})
	t.Run("UnsupportedTypeBool", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := false
		got, err := findNodesLength(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: reflect.ValueOf(val)}, got, err)
	})
	t.Run("SliceType", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		got, err := findNodesLength(selector, reflect.ValueOf([]interface{}{"x", "y"}))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf(2)}, nil, got, err)
	})
	t.Run("MapType", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		got, err := findNodesLength(selector, reflect.ValueOf(map[string]interface{}{
			"x": 1,
			"y": 2,
		}))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf(2)}, nil, got, err)
	})
	t.Run("StringType", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		got, err := findNodesLength(selector, reflect.ValueOf("hello"))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf(5)}, nil, got, err)
	})
}

func TestFindNodesType(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		got, err := findNodesType(selector, nilValue())
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: ".[#]"}, got, err)
	})
	t.Run("Int", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := 0
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("int")}, nil, got, err)
	})
	t.Run("Float", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := 1.1
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("float")}, nil, got, err)
	})
	t.Run("Bool", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := true
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("bool")}, nil, got, err)
	})
	t.Run("String", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := "a"
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("string")}, nil, got, err)
	})
	t.Run("Map", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := map[string]interface{}{"x": 1}
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("map")}, nil, got, err)
	})
	t.Run("Array", func(t *testing.T) {
		selector := Selector{Current: ".[#]", Raw: ".[#]"}
		val := []interface{}{"x"}
		got, err := findNodesType(selector, reflect.ValueOf(val))
		assertQueryMultipleResult(t, []reflect.Value{reflect.ValueOf("array")}, nil, got, err)
	})
}

func TestFindNodesPropertyKeys(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: ".", Raw: "."}
		got, err := findNodesPropertyKeys(selector, nilValue(), false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("UnsupportedCreate", func(t *testing.T) {
		selector := Selector{Current: ".", Raw: "."}
		got, err := findNodesPropertyKeys(selector, reflect.ValueOf(map[string]interface{}{}), true)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedSelector{Selector: selector.Raw}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		previousValue := reflect.ValueOf(0)
		selector := Selector{Current: "x"}
		got, err := findNodesPropertyKeys(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue}, got, err)
	})
	t.Run("SliceValue", func(t *testing.T) {
		previousValue := reflect.ValueOf([]interface{}{"a", "b", "c"})
		selector := Selector{Current: "-"}
		got, err := findNodesPropertyKeys(selector, previousValue, false)

		assertQueryMultipleResult(t, []reflect.Value{
			reflect.ValueOf("0"),
			reflect.ValueOf("1"),
			reflect.ValueOf("2"),
		}, nil, got, err)
	})
	t.Run("MapValue", func(t *testing.T) {
		previousValue := reflect.ValueOf(map[string]interface{}{"name": "Tom", "age": 27})
		selector := Selector{Current: "-"}
		got, err := findNodesPropertyKeys(selector, previousValue, false)

		assertQueryMultipleResultOneOf(t, [][]reflect.Value{
			{
				reflect.ValueOf("name"),
				reflect.ValueOf("age"),
			},
			{
				reflect.ValueOf("age"),
				reflect.ValueOf("name"),
			},
		}, nil, got, err)
	})
}

func TestFindNodesIndex(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: ".", Raw: "."}
		got, err := findNodesIndex(selector, nilValue(), false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		selector := Selector{Current: "[0]", Index: 0, Raw: ".[0]"}
		previousValue := reflect.ValueOf([]interface{}{})
		got, err := findNodesIndex(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: "[0]", PreviousValue: previousValue}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		selector := Selector{Current: "[0]", Index: 0, Raw: ".[0]"}
		previousValue := reflect.ValueOf(map[string]interface{}{})
		got, err := findNodesIndex(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue}, got, err)
	})
}

func TestFindNodesAnyIndex(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		selector := Selector{Current: "[*]", Raw: ".[*]"}
		got, err := findNodesAnyIndex(selector, nilValue())
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: ".[*]"}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		selector := Selector{Current: "[*]", Raw: ".[*]"}
		previousValue := reflect.ValueOf([]interface{}{})
		got, err := findNodesAnyIndex(selector, previousValue)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: "[*]", PreviousValue: previousValue}, got, err)
	})
	t.Run("NotFoundMap", func(t *testing.T) {
		selector := Selector{Current: "[*]", Raw: ".[*]"}
		previousValue := reflect.ValueOf(map[string]interface{}{})
		got, err := findNodesAnyIndex(selector, previousValue)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: "[*]", PreviousValue: previousValue}, got, err)
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		selector := Selector{Current: "[*]", Raw: ".[*]"}
		previousValue := reflect.ValueOf(0)
		got, err := findNodesAnyIndex(selector, previousValue)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue}, got, err)
	})
}

func TestFindNextAvailableIndexNodes(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		previousValue := reflect.ValueOf([]interface{}{})
		selector := Selector{Current: "[0]", Index: 0}
		got, err := findNextAvailableIndexNodes(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}, got, err)
	})
}

func TestFindNodesDynamic(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		previousValue := reflect.ValueOf(nil)
		selector := Selector{Raw: "."}
		got, err := findNodesDynamic(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		previousValue := reflect.ValueOf([]interface{}{})
		selector := Selector{
			Current: "(name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesDynamic(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}, got, err)
	})
	t.Run("NotFoundMap", func(t *testing.T) {
		previousValue := reflect.ValueOf(map[string]interface{}{})
		selector := Selector{
			Current: "(name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesDynamic(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}, got, err)
	})
	t.Run("NotFoundWithCreate", func(t *testing.T) {
		previousValue := reflect.ValueOf([]interface{}{})
		selector := Selector{
			Type:    "NEXT_AVAILABLE_INDEX",
			Current: "(name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesDynamic(selector, previousValue, true)
		if !assertQueryMultipleResult(t, []reflect.Value{nilValue()}, nil, got, err) {
			return
		}
		if exp, got := "NEXT_AVAILABLE_INDEX", selector.Type; exp != got {
			t.Errorf("expected type of %s, got %s", exp, got)
			return
		}
	})
	t.Run("UnsupportedType", func(t *testing.T) {
		previousValue := reflect.ValueOf(0)
		selector := Selector{
			Current: "(name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesDynamic(selector, previousValue, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue}, got, err)
	})
}

func TestFindNodesSearch(t *testing.T) {
	t.Run("NilValue", func(t *testing.T) {
		previousNode := New(nil)
		selector := Selector{Raw: "."}
		got, err := findNodesSearch(selector, previousNode, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnexpectedPreviousNilValue{Selector: "."}, got, err)
	})
	t.Run("NestedNilValue", func(t *testing.T) {
		previousNode := New(map[string]interface{}{
			"x": nil,
		})
		selector := Selector{
			Type:    "SEARCH",
			Raw:     ".(?:name=x)",
			Current: ".(?:name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesSearch(selector, previousNode, false)
		assertQueryMultipleResult(t, []reflect.Value{}, fmt.Errorf("could not find nodes search recursive: %w", &UnexpectedPreviousNilValue{Selector: selector.Current}), got, err)
	})
	t.Run("NotFound", func(t *testing.T) {
		previousNode := New(map[string]interface{}{})
		selector := Selector{
			Type:    "SEARCH",
			Raw:     ".(?:name=x)",
			Current: ".(?:name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesSearch(selector, previousNode, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousNode.Value}, got, err)
	})
	t.Run("NotFoundOptional", func(t *testing.T) {
		previousNode := New(map[string]interface{}{})
		selector := Selector{
			Type:    "SEARCH_OPTIONAL",
			Raw:     ".(#:name=x)",
			Current: ".(#:name=x)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "x"},
			},
		}
		got, err := findNodesSearch(selector, previousNode, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &ValueNotFound{Selector: selector.Current, PreviousValue: previousNode.Value}, got, err)
	})
	t.Run("FoundAllMatches", func(t *testing.T) {
		users := []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "Tom",
			},
			map[string]interface{}{
				"id":   2,
				"name": "Jim",
			},
			map[string]interface{}{
				"id":   3,
				"name": "Tom",
			},
		}
		previousNode := New(users)
		selector := Selector{
			Type:    "SEARCH",
			Raw:     ".(?:name=Tom).name",
			Current: ".(?:name=Tom)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "Tom"},
			},
		}
		got, err := findNodesSearch(selector, previousNode, false)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := []interface{}{
			users[0], users[2],
		}
		gotVals := make([]interface{}, len(got))
		for i, val := range got {
			gotVals[i] = val.Value.Interface()
		}
		if !reflect.DeepEqual(exp, gotVals) {
			t.Errorf("expected %v, got %v", exp, gotVals)
		}
	})
	t.Run("FoundAllMatchesOptional", func(t *testing.T) {
		users := []interface{}{
			map[string]interface{}{
				"id":   1,
				"name": "Tom",
			},
			map[string]interface{}{
				"id":   2,
				"name": "Jim",
			},
			map[string]interface{}{
				"id":   3,
				"name": "Tom",
			},
		}
		previousNode := New(users)
		selector := Selector{
			Type:    "SEARCH_OPTIONAL",
			Raw:     ".(#:name=Tom).name",
			Current: ".(#:name=Tom)",
			Conditions: []Condition{
				&EqualCondition{Key: "name", Value: "Tom"},
			},
		}
		got, err := findNodesSearch(selector, previousNode, false)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		exp := []interface{}{
			users[0], users[2],
		}
		gotVals := make([]interface{}, len(got))
		for i, val := range got {
			gotVals[i] = val.Value.Interface()
		}
		if !reflect.DeepEqual(exp, gotVals) {
			t.Errorf("expected %v, got %v", exp, gotVals)
		}
	})
}

func TestFindNodes(t *testing.T) {
	t.Run("UnsupportedSelector", func(t *testing.T) {
		previousNode := &Node{Value: reflect.ValueOf([]interface{}{})}
		selector := Selector{Raw: "BAD"}
		got, err := findNodes(selector, previousNode, false)
		assertQueryMultipleResult(t, []reflect.Value{}, &UnsupportedSelector{Selector: "BAD"}, got, err)
	})
}
