package model

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
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

func (v Values) ToSliceValue() (*Value, error) {
	slice := NewSliceValue()
	for _, val := range v {
		if err := slice.Append(val); err != nil {
			return nil, err
		}
	}
	return slice, nil
}

// Value represents a value.
type Value struct {
	Value    reflect.Value
	Metadata map[string]any

	setFn func(*Value) error
}

func (v *Value) String() string {
	return v.string(0)
}

func indentStr(indent int) string {
	return strings.Repeat("    ", indent)
}

func (v *Value) string(indent int) string {
	switch v.Type() {
	case TypeString:
		val, err := v.StringValue()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("string{%s}", val)
	case TypeInt:
		val, err := v.IntValue()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("int{%d}", val)
	case TypeFloat:
		val, err := v.FloatValue()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("float(%g)", val)
	case TypeBool:
		val, err := v.BoolValue()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("bool{%t}", val)
	case TypeMap:
		res := fmt.Sprintf("{\n")
		if err := v.RangeMap(func(k string, v *Value) error {
			res += fmt.Sprintf("%s%s: %s,\n", indentStr(indent+1), k, v.string(indent+1))
			return nil
		}); err != nil {
			panic(err)
		}
		return res + indentStr(indent) + "}"
	case TypeSlice:
		md := ""
		if v.IsSpread() {
			md = "spread, "
		}
		if v.IsBranch() {
			md += "branch, "
		}
		res := fmt.Sprintf("array[%s]{\n", strings.TrimSuffix(md, ", "))
		if err := v.RangeSlice(func(k int, v *Value) error {
			res += fmt.Sprintf("%s%d: %s,\n", indentStr(indent+1), k, v.string(indent+1))
			return nil
		}); err != nil {
			panic(err)
		}
		return res + indentStr(indent) + "}"
	case TypeNull:
		return indentStr(indent) + "null"
	default:
		return fmt.Sprintf("unknown[%s]", v.Interface())
	}
}

// NewValue creates a new value.
func NewValue(v any) *Value {
	switch val := v.(type) {
	case *Value:
		return val
	case reflect.Value:
		return &Value{
			Value:    val,
			Metadata: make(map[string]any),
		}
	case nil:
		return NewNullValue()
	default:
		res := newPtr()
		if v != nil {
			res.Elem().Set(reflect.ValueOf(v))
		}
		return &Value{
			Value:    res,
			Metadata: make(map[string]any),
		}
	}
}

// Interface returns the value as an interface.
func (v *Value) Interface() any {
	if v.IsNull() {
		return nil
	}
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
