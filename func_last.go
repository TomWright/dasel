package dasel

import (
	"fmt"
	"reflect"
)

var LastFunc = BasicFunction{
	name: "last",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("last", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				index := val.Len() - 1
				if val.Len() == 0 {
					return nil, fmt.Errorf("index out of range: %w", &ErrIndexNotFound{Index: index})
				}
				value := val.Index(index)
				res = append(res, value)
			default:
				return nil, fmt.Errorf("cannot use last selector on non slice/array types: %w", &ErrIndexNotFound{Index: 0})
			}
		}

		return res, nil
	},
}
