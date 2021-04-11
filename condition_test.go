package dasel_test

import (
	"github.com/tomwright/dasel"
	"reflect"
	"testing"
)

func conditionTest(c dasel.Condition, input interface{}, exp bool, expErr error) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := c.Check(reflect.ValueOf(input))
		if expErr == nil && err != nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err == nil {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if expErr != nil && err != nil && err.Error() != expErr.Error() {
			t.Errorf("expected err %v, got %v", expErr, err)
			return
		}
		if exp != got {
			t.Errorf("expected result %v, got %v", exp, got)
		}
	}
}

func TestEqualCondition_Check(t *testing.T) {
	c := &dasel.EqualCondition{Key: "name", Value: "Tom"}

	t.Run("MatchStringValue", conditionTest(
		&dasel.EqualCondition{Key: "value", Value: "Tom"},
		"Tom",
		true, nil,
	))
	t.Run("MatchMapStringInterface", conditionTest(
		c,
		map[string]interface{}{"name": "Tom"},
		true, nil,
	))
	t.Run("MatchMapInterfaceInterface", conditionTest(
		c,
		map[interface{}]interface{}{"name": "Tom"},
		true, nil,
	))

	t.Run("NoMatchMissingKey", conditionTest(
		c,
		map[string]interface{}{},
		false, nil,
	))
	t.Run("NoMatchMapStringInterface", conditionTest(
		c,
		map[string]interface{}{"name": "Wrong"},
		false, nil,
	))
	t.Run("NoMatchMapInterfaceInterface", conditionTest(
		c,
		map[interface{}]interface{}{"name": "Wrong"},
		false, nil,
	))

	t.Run("Nil", conditionTest(
		c,
		nil,
		false, &dasel.UnhandledCheckType{Value: nil},
	))
	t.Run("String", conditionTest(
		c,
		"",
		false, &dasel.UnhandledCheckType{Value: ""},
	))
}

func TestSortedComparisonCondition_Check(t *testing.T) {
	t.Run("IntNotLessThan", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "3"},
		map[string]interface{}{
			"x": 3,
		},
		false, nil,
	))
	t.Run("IntNotLessThan", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "2"},
		map[string]interface{}{
			"x": 3,
		},
		false, nil,
	))
	t.Run("IntLessThanEqual", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4", Equal: true},
		map[string]interface{}{
			"x": 3,
		},
		true, nil,
	))
	t.Run("IntLessThanEqual", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "3", Equal: true},
		map[string]interface{}{
			"x": 3,
		},
		true, nil,
	))
	t.Run("IntNotMoreThan", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "3", After: true},
		map[string]interface{}{
			"x": 3,
		},
		false, nil,
	))
	t.Run("IntNotMoreThan", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4", After: true},
		map[string]interface{}{
			"x": 3,
		},
		false, nil,
	))
	t.Run("IntMoreThanEqual", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "2", Equal: true, After: true},
		map[string]interface{}{
			"x": 3,
		},
		true, nil,
	))
	t.Run("IntMoreThanEqual", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "3", Equal: true, After: true},
		map[string]interface{}{
			"x": 3,
		},
		true, nil,
	))
	t.Run("MatchMapStringInterfaceLength", conditionTest(
		&dasel.SortedComparisonCondition{Key: "name.[#]", Value: "4"},
		map[string]interface{}{"name": "Tom"},
		true, nil,
	))
	t.Run("MatchMapInterfaceInterfaceLength", conditionTest(
		&dasel.SortedComparisonCondition{Key: "name.[#]", Value: "4"},
		map[interface{}]interface{}{"name": "Tom"},
		true, nil,
	))

	t.Run("NoMatchMissingKey", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4"},
		map[string]interface{}{},
		false, nil,
	))
	t.Run("NoMatchMapStringInterface", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4"},
		map[string]interface{}{"name": "Wrong"},
		false, nil,
	))
	t.Run("NoMatchMapInterfaceInterface", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4"},
		map[interface{}]interface{}{"name": "Wrong"},
		false, nil,
	))

	t.Run("Nil", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4"},
		nil,
		false, &dasel.UnhandledCheckType{Value: nil},
	))
	t.Run("String", conditionTest(
		&dasel.SortedComparisonCondition{Key: "x", Value: "4"},
		"",
		false, &dasel.UnhandledCheckType{Value: ""},
	))
}

func TestKeyEqualCondition_Check(t *testing.T) {
	c := &dasel.KeyEqualCondition{Value: "name"}

	t.Run("MatchStringValue", conditionTest(
		c,
		"name",
		true, nil,
	))
	t.Run("NoMatchMissingKey", conditionTest(
		c,
		"asd",
		false, nil,
	))
	t.Run("Nil", conditionTest(
		c,
		nil,
		false, &dasel.UnhandledCheckType{Value: nil},
	))
}
