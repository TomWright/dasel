package dasel

import (
	"reflect"
)

// KeyEqualCondition lets you check for an exact match.
type KeyEqualCondition struct {
	// Value is the value we are looking for.
	Value string
}

// Check checks to see if other contains the required key value pair.
func (c KeyEqualCondition) Check(other reflect.Value) (bool, error) {
	if !other.IsValid() {
		return false, &UnhandledCheckType{Value: nil}
	}

	value := unwrapValue(other)

	return c.Value == value.String(), nil
}
