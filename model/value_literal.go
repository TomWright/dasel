package model

import (
	"fmt"
	"reflect"
	"slices"
)

func NewNullValue() *Value {
	return NewValue(reflect.New(reflect.TypeFor[any]()))
}

func (v *Value) IsNull() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isNull()
}

func (v *Value) isNull() bool {
	return v.Value.IsNil()
}

func NewStringValue(x string) *Value {
	res := reflect.New(reflect.TypeFor[string]())
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

func (v *Value) IsString() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isString()
}

func (v *Value) isString() bool {
	return v.Value.Kind() == reflect.String
}

func (v *Value) StringValue() (string, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.isString() {
		return "", fmt.Errorf("expected string, got %s", unpacked.Type())
	}
	return unpacked.Value.String(), nil
}

func NewIntValue(x int64) *Value {
	res := reflect.New(reflect.TypeFor[int64]())
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

func (v *Value) IsInt() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isInt()
}

func (v *Value) isInt() bool {
	return slices.Contains([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}, v.Value.Kind())
}

func (v *Value) IntValue() (int64, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.isInt() {
		return 0, fmt.Errorf("expected int, got %s", unpacked.Type())
	}
	return unpacked.Value.Int(), nil
}

func NewFloatValue(x float64) *Value {
	res := reflect.New(reflect.TypeFor[float64]())
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

func (v *Value) IsFloat() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isFloat()
}

func (v *Value) isFloat() bool {
	return slices.Contains([]reflect.Kind{reflect.Float32, reflect.Float64}, v.Value.Kind())
}

func (v *Value) FloatValue() (float64, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.IsFloat() {
		return 0, fmt.Errorf("expected float, got %s", unpacked.Type())
	}
	return unpacked.Value.Float(), nil
}

func NewBoolValue(x bool) *Value {
	res := reflect.New(reflect.TypeFor[bool]())
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

func (v *Value) IsBool() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isBool()
}

func (v *Value) isBool() bool {
	return v.Value.Kind() == reflect.Bool
}

func (v *Value) BoolValue() (bool, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.IsBool() {
		return false, fmt.Errorf("expected bool, got %s", unpacked.Type())
	}
	return unpacked.Value.Bool(), nil
}
