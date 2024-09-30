package model

import (
	"fmt"
	"reflect"
)

func NewSliceValue() *Value {
	s := reflect.MakeSlice(reflect.SliceOf(reflect.TypeFor[any]()), 0, 0)
	ptr := reflect.New(reflect.SliceOf(reflect.TypeFor[any]()))
	ptr.Elem().Set(s)
	return NewValue(ptr)
}

func (v *Value) SliceValue() ([]any, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.IsSlice() {
		return nil, fmt.Errorf("expected slice, got %s", v.Type())
	}
	res, ok := unpacked.Interface().([]any)
	if !ok {
		return nil, fmt.Errorf("could not convert slice to []interface{}")
	}
	return res, nil
}

func (v *Value) IsSlice() bool {
	return v.UnpackKinds(reflect.Interface, reflect.Ptr).isSlice()
}

func (v *Value) isSlice() bool {
	return v.Value.Kind() == reflect.Slice
}

func (v *Value) Append(val *Value) error {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return fmt.Errorf("expected slice, got %s", v.Type())
	}
	newVal := reflect.Append(unpacked.Value, val.Value)
	unpacked.Value.Set(newVal)
	return nil
}

func (v *Value) SliceLen() (int, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return 0, fmt.Errorf("expected slice, got %s", v.Type())
	}
	return unpacked.Value.Len(), nil
}

func (v *Value) GetSliceIndex(i int) (*Value, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return nil, fmt.Errorf("expected slice, got %s", v.Type())
	}
	if i < 0 || i >= unpacked.Value.Len() {
		return nil, fmt.Errorf("index out of range: %d", i)
	}
	return NewValue(unpacked.Value.Index(i)), nil
}
