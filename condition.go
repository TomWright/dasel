package dasel

import "fmt"

type EqualCondition struct {
	Key   string
	Value string
}

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

type Condition interface {
	Check(other interface{}) (bool, error)
}
