package dasel

import (
	"fmt"
	"reflect"
)

// EqualCondition lets you check for an exact match.
type EqualCondition struct {
	// Key is the key of the value to check against.
	Key string
	// Value is the value we are looking for.
	Value string
}

// Check checks to see if other contains the required key value pair.
func (c EqualCondition) Check(other reflect.Value) (bool, error) {
	if !other.IsValid() {
		return false, &UnhandledCheckType{Value: nil}
	}

	value := other
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			if fmt.Sprint(key.Interface()) == c.Key {
				return fmt.Sprint(value.MapIndex(key).Interface()) == c.Value, nil
			}
		}
		return false, nil
	}

	return false, &UnhandledCheckType{Value: value.Kind().String()}
}

// Condition defines a Check we can use within dynamic selectors.
type Condition interface {
	Check(other reflect.Value) (bool, error)
}
