package model

import (
	"reflect"
)

func (v *Value) Set(newValue *Value) error {
	a := v.UnpackKinds(reflect.Ptr, reflect.Interface)
	b := newValue.UnpackKinds(reflect.Ptr, reflect.Interface)

	if a.Kind() == b.Kind() {
		a.Value.Set(b.Value)
		return nil
	}

	x := newPtr()
	x.Elem().Set(b.Value)
	v.Value.Set(x)

	return nil
}
