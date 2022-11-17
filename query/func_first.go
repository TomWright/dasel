package query

import (
	"fmt"
	"reflect"
)

var FirstFunc = BasicFunction{
	name: "first",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected first args given")
		}

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
				return nil, fmt.Errorf("cannot use first selector on non slice/array types")
			}
		}

		return res, nil
	},
}
