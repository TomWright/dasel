package model_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/model"
)

type setTestCase struct {
	valueFn    func() *model.Value
	value      *model.Value
	newValueFn func() *model.Value
	newValue   *model.Value
}

func (tc setTestCase) run(t *testing.T) {
	val := tc.value
	if tc.valueFn != nil {
		val = tc.valueFn()
	}
	newVal := tc.newValue
	if tc.newValueFn != nil {
		newVal = tc.newValueFn()
	}
	if err := val.Set(newVal); err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	eq, err := val.EqualTypeValue(newVal)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if !eq {
		t.Errorf("expected values to be equal")
	}
}

func TestValue_Set(t *testing.T) {
	testCases := []struct {
		name        string
		stringValue func() *model.Value
		intValue    func() *model.Value
		floatValue  func() *model.Value
		boolValue   func() *model.Value
		mapValue    func() *model.Value
		sliceValue  func() *model.Value
	}{
		{
			name: "model constructor",
			stringValue: func() *model.Value {
				return model.NewStringValue("hello")
			},
			intValue: func() *model.Value {
				return model.NewIntValue(1)
			},
			floatValue: func() *model.Value {
				return model.NewFloatValue(1)
			},
			boolValue: func() *model.Value {
				return model.NewBoolValue(true)
			},
			mapValue: func() *model.Value {
				res := model.NewMapValue()
				if err := res.SetMapKey("greeting", model.NewStringValue("hello")); err != nil {
					t.Fatal(err)
				}
				return res
			},
			sliceValue: func() *model.Value {
				res := model.NewSliceValue()
				if err := res.Append(model.NewStringValue("hello")); err != nil {
					t.Fatal(err)
				}
				return res
			},
		},
		{
			name: "go types non ptr",
			stringValue: func() *model.Value {
				v := "hello"
				return model.NewValue(v)
			},
			intValue: func() *model.Value {
				v := int64(1)
				return model.NewValue(v)
			},
			floatValue: func() *model.Value {
				v := 1.0
				return model.NewValue(v)
			},
			boolValue: func() *model.Value {
				v := true
				return model.NewValue(v)
			},
			mapValue: func() *model.Value {
				v := map[string]interface{}{
					"greeting": "hello",
				}
				return model.NewValue(v)
			},
			sliceValue: func() *model.Value {
				v := []interface{}{
					"hello",
				}
				return model.NewValue(v)
			},
		},
		{
			name: "go types ptr",
			stringValue: func() *model.Value {
				v := "hello"
				return model.NewValue(&v)
			},
			intValue: func() *model.Value {
				v := int64(1)
				return model.NewValue(&v)
			},
			floatValue: func() *model.Value {
				v := 1.0
				return model.NewValue(&v)
			},
			boolValue: func() *model.Value {
				v := true
				return model.NewValue(&v)
			},
			mapValue: func() *model.Value {
				v := map[string]interface{}{
					"greeting": "hello",
				}
				return model.NewValue(&v)
			},
			sliceValue: func() *model.Value {
				v := []interface{}{
					"hello",
				}
				return model.NewValue(&v)
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Run("string", setTestCase{
				valueFn:  tc.stringValue,
				newValue: model.NewStringValue("world"),
			}.run)
			t.Run("int", setTestCase{
				valueFn:  tc.intValue,
				newValue: model.NewIntValue(2),
			}.run)
			t.Run("float", setTestCase{
				valueFn:  tc.floatValue,
				newValue: model.NewFloatValue(2),
			}.run)
			t.Run("bool", setTestCase{
				valueFn:  tc.boolValue,
				newValue: model.NewBoolValue(false),
			}.run)
			t.Run("map", setTestCase{
				valueFn: tc.mapValue,
				newValueFn: func() *model.Value {
					res := model.NewMapValue()
					if err := res.SetMapKey("greeting", model.NewStringValue("world")); err != nil {
						t.Fatal(err)
					}
					return res
				},
			}.run)
			t.Run("slice", setTestCase{
				valueFn: tc.sliceValue,
				newValueFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewStringValue("world")); err != nil {
						t.Fatal(err)
					}
					return res
				},
			}.run)
			t.Run("string over int", setTestCase{
				valueFn:  tc.intValue,
				newValue: model.NewStringValue("world"),
			}.run)
			t.Run("int over float", setTestCase{
				valueFn:  tc.floatValue,
				newValue: model.NewIntValue(2),
			}.run)
			t.Run("float over bool", setTestCase{
				valueFn:  tc.boolValue,
				newValue: model.NewFloatValue(2),
			}.run)
			t.Run("bool over map", setTestCase{
				valueFn:  tc.mapValue,
				newValue: model.NewBoolValue(true),
			}.run)
			t.Run("map over slice", setTestCase{
				valueFn: tc.sliceValue,
				newValueFn: func() *model.Value {
					res := model.NewMapValue()
					if err := res.SetMapKey("greeting", model.NewStringValue("world")); err != nil {
						t.Fatal(err)
					}
					return res
				},
			}.run)
			t.Run("string over slice", setTestCase{
				valueFn:  tc.sliceValue,
				newValue: model.NewStringValue("world"),
			}.run)
			t.Run("slice over map", setTestCase{
				valueFn: tc.mapValue,
				newValueFn: func() *model.Value {
					res := model.NewSliceValue()
					if err := res.Append(model.NewStringValue("world")); err != nil {
						t.Fatal(err)
					}
					return res
				},
			}.run)
		})
	}
}
