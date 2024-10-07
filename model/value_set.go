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

	// todo : figure this out
	x := newPtr()
	x.Elem().Set(b.Value)

	target, err := v.UnpackUntilAddressable()
	if err != nil {
		return err
	}

	target.Value.Set(x)

	return nil
}
