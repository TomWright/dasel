package model

import (
	"fmt"
	"reflect"
)

// Set sets the value of the value.
func (v *Value) Set(newValue *Value) error {
	if v.setFn != nil {
		return v.setFn(newValue)
	}

	a, err := v.UnpackUntilAddressable()
	if err != nil {
		return err
	}

	b := newValue.UnpackKinds(reflect.Ptr)
	if a.Kind() == b.Kind() {
		a.Value.Set(b.Value)
		return nil
	}

	b = newValue.UnpackKinds(reflect.Ptr, reflect.Interface)
	if a.Kind() == b.Kind() {
		a.Value.Set(b.Value)
		return nil
	}

	// These are commented out because I don't think they are needed.

	//if a.Kind() == newValue.Kind() {
	//	a.Value.Set(newValue.Value)
	//	return nil
	//}

	//b = newValue.UnpackKinds(reflect.Interface)
	//if a.Kind() == b.Kind() {
	//	a.Value.Set(b.Value)
	//	return nil
	//}

	//b = newValue.UnpackKinds(reflect.Ptr, reflect.Interface)
	//if a.Kind() == b.Kind() {
	//	a.Value.Set(b.Value)
	//	return nil
	//}

	//b, err = newValue.UnpackUntilAddressable()
	//if err != nil {
	//	return err
	//}
	//if a.Kind() == b.Kind() {
	//	a.Value.Set(b.Value)
	//	return nil
	//}

	// This is a hard limitation at the moment.
	// If the types are not the same, we cannot set the value.
	return fmt.Errorf("could not set %s value on %s value", newValue.Type(), v.Type())
}
