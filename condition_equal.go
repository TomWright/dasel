package dasel

import (
	"errors"
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

	value := unwrapValue(other)

	if c.Key == "value" || c.Key == "." {
		return fmt.Sprint(value.Interface()) == c.Value, nil
	}

	switch value.Kind() {
	case reflect.Map, reflect.Slice:
		subRootNode := New(value.Interface())
		foundNode, err := subRootNode.Query(c.Key)
		if err != nil {
			var valueNotFound = &ValueNotFound{}
			if errors.As(err, &valueNotFound) {
				return false, nil
			}

			return false, fmt.Errorf("subquery failed: %w", err)
		}

		return fmt.Sprint(foundNode.InterfaceValue()) == c.Value, nil
	}

	return false, &UnhandledCheckType{Value: value.Kind().String()}
}
