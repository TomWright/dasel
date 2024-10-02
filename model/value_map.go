package model

import (
	"fmt"
	"reflect"

	"github.com/tomwright/dasel/v3/dencoding"
)

func NewMapValue() *Value {
	return NewValue(dencoding.NewMap())
}

func (v *Value) IsMap() bool {
	return v.isStandardMap() || v.isDencodingMap()
}

func (v *Value) isStandardMap() bool {
	return v.UnpackKinds(reflect.Interface, reflect.Ptr).Kind() == reflect.Map
}

func (v *Value) isDencodingMap() bool {
	return v.UnpackKinds(reflect.Interface, reflect.Ptr).Value.Type() == reflect.TypeFor[dencoding.Map]()
}

func (v *Value) dencodingMapValue() (*dencoding.Map, error) {
	if v.isDencodingMap() {
		m, err := v.UnpackUntilType(reflect.TypeFor[*dencoding.Map]())
		if err != nil {
			return nil, fmt.Errorf("error getting map: %w", err)
		}
		return m.Value.Interface().(*dencoding.Map), nil
	}
	return nil, fmt.Errorf("value is not a dencoding map")
}

// SetMapKey sets the value at the specified key in the map.
func (v *Value) SetMapKey(key string, value *Value) error {
	switch {
	case v.isDencodingMap():
		m, err := v.dencodingMapValue()
		if err != nil {
			return fmt.Errorf("error getting map: %w", err)
		}
		m.Set(key, value.Value.Interface())
		return nil
	case v.isStandardMap():
		unpacked, err := v.UnpackUntilKind(reflect.Map)
		if err != nil {
			return fmt.Errorf("error unpacking value: %w", err)
		}
		unpacked.Value.SetMapIndex(reflect.ValueOf(key), value.Value)
		return nil
	default:
		return fmt.Errorf("value is not a map")
	}
}

// GetMapKey returns the value at the specified key in the map.
func (v *Value) GetMapKey(key string) (*Value, error) {
	switch {
	case v.isDencodingMap():
		m, err := v.dencodingMapValue()
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
	case v.isStandardMap():
		unpacked, err := v.UnpackUntilKind(reflect.Map)
		if err != nil {
			return nil, fmt.Errorf("error unpacking value: %w", err)
		}
		i := unpacked.Value.MapIndex(reflect.ValueOf(key))
		if !i.IsValid() {
			return nil, &MapKeyNotFound{Key: key}
		}
		return &Value{
			Value: i,
		}, nil
	default:
		return nil, fmt.Errorf("value is not a map")
	}
}

// DeleteMapKey deletes the key from the map.
func (v *Value) DeleteMapKey(key string) error {
	switch {
	case v.isDencodingMap():
		m, err := v.dencodingMapValue()
		if err != nil {
			return fmt.Errorf("error getting map: %w", err)
		}
		m.Delete(key)
		return nil
	case v.isStandardMap():
		unpacked, err := v.UnpackUntilKind(reflect.Map)
		if err != nil {
			return fmt.Errorf("error unpacking value: %w", err)
		}
		unpacked.Value.SetMapIndex(reflect.ValueOf(key), reflect.Value{})
		return nil
	default:
		return fmt.Errorf("value is not a map")
	}
}

// MapKeys returns a list of keys in the map.
func (v *Value) MapKeys() ([]string, error) {
	switch {
	case v.isDencodingMap():
		m, err := v.dencodingMapValue()
		if err != nil {
			return nil, fmt.Errorf("error getting map: %w", err)
		}
		return m.Keys(), nil
	case v.isStandardMap():
		unpacked, err := v.UnpackUntilKind(reflect.Map)
		if err != nil {
			return nil, fmt.Errorf("error unpacking value: %w", err)
		}
		keys := unpacked.Value.MapKeys()
		strKeys := make([]string, len(keys))
		for i, k := range keys {
			strKeys[i] = k.String()
		}
		return strKeys, nil
	default:
		return nil, fmt.Errorf("value is not a map")
	}
}

// RangeMap iterates over each key in the map and calls the provided function with the key and value.
func (v *Value) RangeMap(f func(string, *Value) error) error {
	keys, err := v.MapKeys()
	if err != nil {
		return fmt.Errorf("error getting map keys: %w", err)
	}

	for _, k := range keys {
		va, err := v.GetMapKey(k)
		if err != nil {
			return fmt.Errorf("error getting map key: %w", err)
		}
		if err := f(k, va); err != nil {
			return err
		}
	}

	return nil
}

// MapKeyValues returns a list of key value pairs in the map.
func (v *Value) MapKeyValues() ([]KeyValue, error) {
	keys, err := v.MapKeys()
	if err != nil {
		return nil, fmt.Errorf("error getting map keys: %w", err)
	}

	kvs := make([]KeyValue, len(keys))

	for _, k := range keys {
		va, err := v.GetMapKey(k)
		if err != nil {
			return nil, fmt.Errorf("error getting map key: %w", err)
		}
		kvs = append(kvs, KeyValue{
			Key:   k,
			Value: va,
		})
	}

	return kvs, nil
}
