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

// ToSliceValue converts a list of values to a slice value.
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
	value    reflect.Value
	Metadata map[string]any

	setFn func(*Value) error
}

// String returns the value as a formatted string, along with type info.
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
		return NewNestedValue(val)
	case reflect.Value:
		return &Value{
			value:    val,
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
			value:    res,
			Metadata: make(map[string]any),
		}
	}
}

// NewNestedValue creates a new nested value.
func NewNestedValue(v *Value) *Value {
	return &Value{
		value:    reflect.ValueOf(v),
		Metadata: make(map[string]any),
	}
}

func (v *Value) isDaselValue() bool {
	cur := v.value
	for cur.Kind() == reflect.Interface && !cur.IsNil() {
		cur = cur.Elem()
	}
	return cur.Type() == reflect.TypeFor[*Value]()
}

func (v *Value) daselValue() (*Value, error) {
	if v.isDaselValue() {
		m, err := v.UnpackUntilType(reflect.TypeFor[*Value]())
		if err != nil {
			return nil, fmt.Errorf("error getting dasel value: %w", err)
		}
		return m.value.Interface().(*Value), nil
	}
	return nil, fmt.Errorf("value is not a dasel value")
}

// Interface returns the value as an interface.
func (v *Value) Interface() any {
	if v.IsNull() {
		return nil
	}
	return v.value.Interface()
}

// Kind returns the reflect kind of the value.
func (v *Value) Kind() reflect.Kind {
	return v.value.Kind()
}

// UnpackKinds unpacks the reflect value until it no longer matches the given kinds.
func (v *Value) UnpackKinds(kinds ...reflect.Kind) *Value {
	val := v
	for val.isDaselValue() {
		var err error
		val, err = val.daselValue()
		if err != nil {
			panic(err)
		}
	}
	res := val.value
	for {
		if !slices.Contains(kinds, res.Kind()) || res.IsNil() {
			return NewValue(res)
		}
		res = res.Elem()
	}
}

type ErrCouldNotUnpackToType struct {
	Type reflect.Type
}

func (e ErrCouldNotUnpackToType) Error() string {
	return fmt.Sprintf("could not unpack to type: %s", e.Type)
}

// UnpackUntilType unpacks the reflect value until it matches the given type.
func (v *Value) UnpackUntilType(t reflect.Type) (*Value, error) {
	res := v.value
	for {
		if res.Type() == t {
			return NewValue(res), nil
		}
		if res.Kind() == reflect.Interface || res.Kind() == reflect.Ptr && !res.IsNil() {
			res = res.Elem()
			continue
		}
		return nil, &ErrCouldNotUnpackToType{Type: t}
	}
}

// UnpackUntilAddressable unpacks the reflect value until it is addressable.
func (v *Value) UnpackUntilAddressable() (*Value, error) {
	res := v.value
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
	res := v.value
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
	res := v.value
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

// IsScalar returns true if the type is scalar.
func (v *Value) IsScalar() bool {
	switch {
	case v.IsString():
		return true
	case v.IsInt():
		return true
	case v.IsFloat():
		return true
	case v.IsBool():
		return true
	case v.IsNull():
		return true
	default:
		return false
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

func (v *Value) Copy() (*Value, error) {
	switch v.Type() {
	case TypeMap:
		return v.MapCopy()
	default:
		return nil, fmt.Errorf("copy not supported for type: %s", v.Type())
	}
}
