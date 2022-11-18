package dasel

import (
	"fmt"
	"reflect"
)

var AllFunc = BasicFunction{
	name: "all",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected last args given")
		}

		res := make(Values, 0)

		for _, val := range input {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i < val.Len(); i++ {
					res = append(res, val.Index(i))
				}
			case reflect.Map:
				for _, key := range val.MapKeys() {
					res = append(res, val.MapIndex(key))
				}
			default:
				return nil, fmt.Errorf("cannot use all selector on non slice/array/map types")
			}
		}

		return res, nil
	},
}
