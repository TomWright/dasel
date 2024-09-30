package model

import (
	"reflect"
	"slices"
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
	TypeMap     Type = "map"
	TypeSlice   Type = "array"
	TypeUnknown Type = "unknown"
)

type Value struct {
	Value reflect.Value
}

func NewValue(v any) *Value {
	if rv, ok := v.(reflect.Value); ok {
		return &Value{
			Value: rv,
		}
	}
	return &Value{
		Value: reflect.ValueOf(v),
	}
}

func (v *Value) Interface() interface{} {
	return v.Value.Interface()
}

func (v *Value) UnpackKinds(kinds ...reflect.Kind) *Value {
	res := v.Value
	for {
		if !slices.Contains(kinds, res.Kind()) {
			return NewValue(res)
		}
		res = res.Elem()
	}
}

func (v *Value) Type() Type {
	switch {
	case v.IsString():
		return TypeString
	case v.IsInt():
		return TypeInt
	case v.IsFloat():
		return TypeFloat
	case v.IsBool():
		return TypeBool
	case v.IsMap():
		return TypeMap
	case v.IsSlice():
		return TypeSlice
	default:
		return TypeUnknown
	}
}
