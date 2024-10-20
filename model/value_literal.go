package model

import (
	"fmt"
	"reflect"
	"slices"
)

func newPtr() reflect.Value {
	return reflect.New(reflect.TypeFor[any]())
}

// NewNullValue creates a new Value with a nil value.
func NewNullValue() *Value {
	return NewValue(newPtr())
}

// IsNull returns true if the value is null.
func (v *Value) IsNull() bool {
	return v.isNull()
}

func (v *Value) isNull() bool {
	return v.Value.IsNil()
}

// NewStringValue creates a new Value with a string value.
func NewStringValue(x string) *Value {
	res := newPtr()
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

// IsString returns true if the value is a string.
func (v *Value) IsString() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isString()
}

func (v *Value) isString() bool {
	return v.Value.Kind() == reflect.String
}

// StringValue returns the string value of the Value.
func (v *Value) StringValue() (string, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.isString() {
		return "", fmt.Errorf("expected string, got %s", unpacked.Type())
	}
	return unpacked.Value.String(), nil
}

// StringLen returns the length of the string.
func (v *Value) StringLen() (int, error) {
	val, err := v.StringValue()
	if err != nil {
		return 0, err
	}
	return len(val), nil
}

// StringIndexRange returns a new string containing the values between the start and end indexes.
// Comparable to go's string[start:end].
func (v *Value) StringIndexRange(start, end int) (*Value, error) {
	strVal, err := v.StringValue()
	if err != nil {
		return nil, err
	}

	inBytes := []rune(strVal)
	l := len(inBytes)

	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}

	resBytes := make([]rune, 0)

	if start > end {
		for i := start; i >= end; i-- {
			resBytes = append(resBytes, inBytes[i])
		}
	} else {
		for i := start; i <= end; i++ {
			resBytes = append(resBytes, inBytes[i])
		}
	}

	res := string(resBytes)

	return NewStringValue(res), nil
}

// NewIntValue creates a new Value with an int value.
func NewIntValue(x int64) *Value {
	res := newPtr()
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

// IsInt returns true if the value is an int.
func (v *Value) IsInt() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isInt()
}

func (v *Value) isInt() bool {
	return slices.Contains([]reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64}, v.Value.Kind())
}

// IntValue returns the int value of the Value.
func (v *Value) IntValue() (int64, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.isInt() {
		return 0, fmt.Errorf("expected int, got %s", unpacked.Type())
	}
	return unpacked.Value.Int(), nil
}

// NewFloatValue creates a new Value with a float value.
func NewFloatValue(x float64) *Value {
	res := newPtr()
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

// IsFloat returns true if the value is a float.
func (v *Value) IsFloat() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isFloat()
}

func (v *Value) isFloat() bool {
	return slices.Contains([]reflect.Kind{reflect.Float32, reflect.Float64}, v.Value.Kind())
}

// FloatValue returns the float value of the Value.
func (v *Value) FloatValue() (float64, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.IsFloat() {
		return 0, fmt.Errorf("expected float, got %s", unpacked.Type())
	}
	return unpacked.Value.Float(), nil
}

// NewBoolValue creates a new Value with a bool value.
func NewBoolValue(x bool) *Value {
	res := newPtr()
	res.Elem().Set(reflect.ValueOf(x))
	return NewValue(res)
}

// IsBool returns true if the value is a bool.
func (v *Value) IsBool() bool {
	return v.UnpackKinds(reflect.Ptr, reflect.Interface).isBool()
}

func (v *Value) isBool() bool {
	return v.Value.Kind() == reflect.Bool
}

// BoolValue returns the bool value of the Value.
func (v *Value) BoolValue() (bool, error) {
	unpacked := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	if !unpacked.IsBool() {
		return false, fmt.Errorf("expected bool, got %s", unpacked.Type())
	}
	return unpacked.Value.Bool(), nil
}
