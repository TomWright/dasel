package dasel

import (
	"fmt"
	"reflect"
)

var FirstFunc = BasicFunction{
	name: "first",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("first", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				if val.Len() == 0 {
					return nil, fmt.Errorf("index out of range: %w", &ErrIndexNotFound{Index: 0})
				}
				value := val.Index(0)
				res = append(res, value)
			default:
				return nil, fmt.Errorf("cannot use first selector on non slice/array types: %w", &ErrIndexNotFound{Index: 0})
			}
		}

		return res, nil
	},
}
