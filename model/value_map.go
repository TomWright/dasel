package model

import (
	"fmt"
	"reflect"

	"github.com/tomwright/dasel/v3/dencoding"
)

func NewMapValue() *Value {
	return NewValue(dencoding.NewMap())
}

func (v *Value) MapValue() (*dencoding.Map, error) {
	if !v.IsMap() {
		return nil, fmt.Errorf("value is not a map")
	}
	return v.Value.Interface().(*dencoding.Map), nil
}

func (v *Value) IsMap() bool {
	return v.UnpackKinds(reflect.Interface).isMap()
}

func (v *Value) isMap() bool {
	return v.Value.Type() == reflect.TypeFor[*dencoding.Map]()
}

func (v *Value) SetMapKey(key string, value *Value) error {
	m, err := v.MapValue()
	if err != nil {
		return fmt.Errorf("error getting map: %w", err)
	}
	m.Set(key, value.Value.Interface())
	return nil
}

func (v *Value) GetMapKey(key string) (*Value, error) {
	m, err := v.MapValue()
	if err != nil {
		return nil, fmt.Errorf("error getting map: %w", err)
	}
	val, ok := m.Get(key)
	if !ok {
		return nil, &MapKeyNotFound{Key: key}
	}
	return &Value{
		Value: reflect.ValueOf(val),
	}, nil
}
