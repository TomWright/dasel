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
	// Not is true if this is a not equal check.
	Not bool
}

func (c EqualCondition) check(a interface{}, b interface{}) (bool, error) {
	var res = fmt.Sprint(a) == b
	if c.Not {
		res = !res
	}
	return res, nil
}

// Check checks to see if other contains the required key value pair.
func (c EqualCondition) Check(other reflect.Value) (bool, error) {
	if !other.IsValid() {
		return false, &UnhandledCheckType{Value: nil}
	}

	value := unwrapValue(other)

	if c.Key == "value" || c.Key == "." {
		return c.check(value.Interface(), c.Value)
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

		return c.check(foundNode.InterfaceValue(), c.Value)
	}

	return false, &UnhandledCheckType{Value: value.String()}
}
