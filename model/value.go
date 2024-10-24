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
	TypeNull    Type = "null"
)

// KeyValue represents a key value pair.
type KeyValue struct {
	Key   string
	Value *Value
}

// Values represents a list of values.
type Values []*Value

// Value represents a value.
type Value struct {
	Value    reflect.Value
	Metadata map[string]any

	setFn func(*Value) error
}

// NewValue creates a new value.
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

// Interface returns the value as an interface.
func (v *Value) Interface() any {
	return v.Value.Interface()
}

// Kind returns the reflect kind of the value.
func (v *Value) Kind() reflect.Kind {
	return v.Value.Kind()
}

// UnpackKinds unpacks the reflect value until it no longer matches the given kinds.
func (v *Value) UnpackKinds(kinds ...reflect.Kind) *Value {
	res := v.Value
	for {
		if !slices.Contains(kinds, res.Kind()) || res.IsNil() {
			return NewValue(res)
		}
		res = res.Elem()
	}
}

// UnpackUntilType unpacks the reflect value until it matches the given type.
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

// UnpackUntilAddressable unpacks the reflect value until it is addressable.
func (v *Value) UnpackUntilAddressable() (*Value, error) {
	res := v.Value
	for {
		if res.CanAddr() {
			return NewValue(res), nil
		}
		if res.Kind() == reflect.Interface || res.Kind() == reflect.Ptr && !res.IsNil() {
			res = res.Elem()
			continue
		}
		return nil, fmt.Errorf("could not unpack addressable value")
	}
}

// UnpackUntilKind unpacks the reflect value until it matches the given kind.
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

// UnpackUntilKinds unpacks the reflect value until it matches the given kind.
func (v *Value) UnpackUntilKinds(kinds ...reflect.Kind) (*Value, error) {
	res := v.Value
	for {
		if slices.Contains(kinds, res.Kind()) {
			return NewValue(res), nil
		}
		if res.Kind() == reflect.Interface || res.Kind() == reflect.Ptr && !res.IsNil() {
			res = res.Elem()
			continue
		}
		return nil, fmt.Errorf("could not unpack to kinds: %v", kinds)
	}
}

// Type returns the type of the value.
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
	case v.IsNull():
		return TypeNull
	default:
		return TypeUnknown
	}
}

// Len returns the length of the value.
func (v *Value) Len() (int, error) {
	var l int
	var err error

	switch {
	case v.IsSlice():
		l, err = v.SliceLen()
	case v.IsMap():
		l, err = v.MapLen()
	case v.IsString():
		l, err = v.StringLen()
	default:
		err = ErrUnexpectedTypes{
			Expected: []Type{TypeSlice, TypeMap, TypeString},
			Actual:   v.Type(),
		}
	}

	if err != nil {
		return l, err
	}

	return l, nil
}
