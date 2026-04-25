package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

func TestTernary(t *testing.T) {
	t.Run("true literal", testCase{
		s:   `true ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("false literal", testCase{
		s:   `false ? "yes" : "no"`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("condition with comparison true", testCase{
		s:   `1 + 1 == 2 ? "math" : "broken"`,
		out: model.NewStringValue("math"),
	}.run)
	t.Run("condition with comparison false", testCase{
		s:   `1 + 1 == 3 ? "math" : "broken"`,
		out: model.NewStringValue("broken"),
	}.run)
	t.Run("nested ternary in then", testCase{
		s:   `true ? (false ? "a" : "b") : "c"`,
		out: model.NewStringValue("b"),
	}.run)
	t.Run("nested ternary in else", testCase{
		s:   `false ? "a" : (true ? "b" : "c")`,
		out: model.NewStringValue("b"),
	}.run)
	t.Run("double nested ternary", testCase{
		s:   `true ? (true ? (false ? "a" : "b") : "c") : "d"`,
		out: model.NewStringValue("b"),
	}.run)
	t.Run("with input data true", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"age": 25})
		},
		s:   `age >= 18 ? "adult" : "minor"`,
		out: model.NewStringValue("adult"),
	}.run)
	t.Run("with input data false", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"age": 10})
		},
		s:   `age >= 18 ? "adult" : "minor"`,
		out: model.NewStringValue("minor"),
	}.run)
	t.Run("returns number", testCase{
		s:   `true ? 42 : 0`,
		out: model.NewIntValue(42),
	}.run)
	t.Run("returns float", testCase{
		s:   `true ? 3.14 : 0.0`,
		out: model.NewFloatValue(3.14),
	}.run)
	t.Run("returns null in else", testCase{
		s:   `false ? "yes" : null`,
		out: model.NewNullValue(),
	}.run)
	t.Run("returns null in then", testCase{
		s:   `true ? null : "no"`,
		out: model.NewNullValue(),
	}.run)
	t.Run("returns property", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"active": true, "name": "Alice"})
		},
		s:   `active ? name : "unknown"`,
		out: model.NewStringValue("Alice"),
	}.run)
	t.Run("returns property from else", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"active": false, "name": "Alice", "fallback": "Bob"})
		},
		s:   `active ? name : fallback`,
		out: model.NewStringValue("Bob"),
	}.run)
	t.Run("chained property in condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{
				"user": map[string]interface{}{"active": true},
			})
		},
		s:   `user.active ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("arithmetic in then branch", testCase{
		s:   `true ? 1 + 2 : 3 + 4`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("arithmetic in else branch", testCase{
		s:   `false ? 1 + 2 : 3 + 4`,
		out: model.NewIntValue(7),
	}.run)
	t.Run("logical and condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"a": 5, "b": 10})
		},
		s:   `a > 1 && b > 5 ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("logical and condition false", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"a": 0, "b": 10})
		},
		s:   `a > 1 && b > 5 ? "yes" : "no"`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("logical or condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"a": 0, "b": 10})
		},
		s:   `a > 1 || b > 5 ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("logical or condition false", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"a": 0, "b": 1})
		},
		s:   `a > 1 || b > 5 ? "yes" : "no"`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("func call in then branch", testCase{
		s:   `true ? len("abc") : len("de")`,
		out: model.NewIntValue(3),
	}.run)
	t.Run("func call in else branch", testCase{
		s:   `false ? len("abc") : len("de")`,
		out: model.NewIntValue(2),
	}.run)
	t.Run("string concatenation in then", testCase{
		s:   `true ? "hello" + " " + "world" : "goodbye"`,
		out: model.NewStringValue("hello world"),
	}.run)
	t.Run("string concatenation in else", testCase{
		s:   `false ? "hello" : "good" + "bye"`,
		out: model.NewStringValue("goodbye"),
	}.run)
	t.Run("equality check false condition", testCase{
		s:   `1 == 2 ? "yes" : "no"`,
		out: model.NewStringValue("no"),
	}.run)
	t.Run("not equal condition", testCase{
		s:   `1 != 2 ? "different" : "same"`,
		out: model.NewStringValue("different"),
	}.run)
	t.Run("less than condition", testCase{
		s:   `1 < 2 ? "less" : "not less"`,
		out: model.NewStringValue("less"),
	}.run)
	t.Run("greater than or equal condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"val": 5})
		},
		s:   `val >= 5 ? "gte" : "lt"`,
		out: model.NewStringValue("gte"),
	}.run)
	t.Run("less than or equal condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"val": 5})
		},
		s:   `val <= 5 ? "lte" : "gt"`,
		out: model.NewStringValue("lte"),
	}.run)
	t.Run("subtraction in condition", testCase{
		s:   `10 - 3 == 7 ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("bool property as condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"enabled": true})
		},
		s:   `enabled ? "on" : "off"`,
		out: model.NewStringValue("on"),
	}.run)
	t.Run("bool property false as condition", testCase{
		inFn: func() *model.Value {
			return model.NewValue(map[string]interface{}{"enabled": false})
		},
		s:   `enabled ? "on" : "off"`,
		out: model.NewStringValue("off"),
	}.run)
	t.Run("multiplication in condition", testCase{
		s:   `2 * 3 == 6 ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
	t.Run("modulo in condition", testCase{
		s:   `10 % 3 == 1 ? "yes" : "no"`,
		out: model.NewStringValue("yes"),
	}.run)
}
