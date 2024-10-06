package model

import (
	"fmt"
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

type KeyValue struct {
	Key   string
	Value *Value
}

type Value struct {
	Value    reflect.Value
	Metadata map[string]any
}

func NewValue(v any) *Value {
	switch val := v.(type) {
	case *Value:
		return val
	case reflect.Value:
		return &Value{
			Value: val,
		}
	case nil:
		return NewNullValue()
	default:
		res := newPtr()
		if v != nil {
			res.Elem().Set(reflect.ValueOf(v))
		}
		return &Value{
			Value: res,
		}
	}
}

func (v *Value) Interface() interface{} {
	return v.Value.Interface()
}

func (v *Value) Kind() reflect.Kind {
	return v.Value.Kind()
}

func (v *Value) UnpackKinds(kinds ...reflect.Kind) *Value {
	res := v.Value
	for {
		if !slices.Contains(kinds, res.Kind()) || res.IsNil() {
			return NewValue(res)
		}
		res = res.Elem()
	}
}

func (v *Value) UnpackUntilType(t reflect.Type) (*Value, error) {
	res := v.Value
	for {
		if res.Type() == t {
			return NewValue(res), nil
		}
		if res.Kind() == reflect.Interface || res.Kind() == reflect.Ptr && !res.IsNil() {
			res = res.Elem()
			continue
		}
		return nil, fmt.Errorf("could not unpack to type: %s", t)
	}
}

func (v *Value) UnpackUntilKind(k reflect.Kind) (*Value, error) {
	res := v.Value
	for {
		if res.Kind() == k {
			return NewValue(res), nil
		}
		if res.Kind() == reflect.Interface || res.Kind() == reflect.Ptr && !res.IsNil() {
			res = res.Elem()
			continue
		}
		return nil, fmt.Errorf("could not unpack to kind: %s", k)
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

func (v *Value) Len() (int, error) {
	var l int
	var err error

	switch {
	case v.IsSlice():
		l, err = v.SliceLen()
	case v.IsMap():
		l, err = v.MapLen()
	default:
		err = fmt.Errorf("len expects slice or map")
	}

	if err != nil {
		return l, err
	}

	return l, nil
}
