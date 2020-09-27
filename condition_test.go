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
