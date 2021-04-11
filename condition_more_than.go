package dasel

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// SortedComparisonCondition lets you check for an exact match.
type SortedComparisonCondition struct {
	// Key is the key of the value to check against.
	Key string
	// Value is the value we are looking for.
	Value string
	// Equal is true if the values can match.
	Equal bool
	// After is true if the input value should be sorted after the Value.
	After bool
}

// Check checks to see if other contains the required key value pair.
func (c SortedComparisonCondition) Check(other reflect.Value) (bool, error) {
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

		foundValueStr := fmt.Sprint(foundNode.InterfaceValue())

		// Check if the values are equal
		if foundValueStr == c.Value {
			return c.Equal, nil
		}

		sortedVals := []string{foundValueStr, c.Value}
		sort.Strings(sortedVals)

		if !c.After && sortedVals[1] == c.Value {
			return true, nil
		} else if c.After && sortedVals[0] == c.Value {
			return true, nil
		}

		return false, nil
	}

	return false, &UnhandledCheckType{Value: value.String()}
}
