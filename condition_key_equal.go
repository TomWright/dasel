package dasel

import (
	"fmt"
	"reflect"
)

// KeyEqualCondition lets you check for an exact match.
type KeyEqualCondition struct {
	// Value is the value we are looking for.
	Value string
	// Not is true if this is a not equal check.
	Not bool
}

func (c KeyEqualCondition) check(a interface{}, b interface{}) (bool, error) {
	var res = fmt.Sprint(a) == b
	if c.Not {
		res = !res
	}
	return res, nil
}

// Check checks to see if other contains the required key value pair.
func (c KeyEqualCondition) Check(other reflect.Value) (bool, error) {
	if !other.IsValid() {
		return false, &UnhandledCheckType{Value: nil}
	}

	value := unwrapValue(other)

	return c.check(c.Value, value.String())
}
