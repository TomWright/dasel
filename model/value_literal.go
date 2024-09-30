package model

import (
	"fmt"
	"reflect"
)

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
	return v.Value.Kind() == reflect.Int64
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
	return v.Value.Kind() == reflect.Float64
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
