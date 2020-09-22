package dasel

import "fmt"

// EqualCondition lets you check for an exact match.
type EqualCondition struct {
	// Key is the key of the value to check against.
	Key string
	// Value is the value we are looking for.
	Value string
}

// Check checks to see if other contains the required key value pair.
func (c EqualCondition) Check(other interface{}) (bool, error) {
	switch o := other.(type) {
	case map[string]string:
		return o[c.Key] == c.Value, nil
	case map[string]interface{}:
		return fmt.Sprint(o[c.Key]) == c.Value, nil
	case map[interface{}]interface{}:
		return fmt.Sprint(o[c.Key]) == c.Value, nil
	default:
		return false, fmt.Errorf("unhandled check type: %T", other)
	}
}

// Condition defines a Check we can use within dynamic selectors.
type Condition interface {
	Check(other interface{}) (bool, error)
}
