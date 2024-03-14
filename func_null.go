package dasel

import (
	"reflect"
)

var NullFunc = BasicFunction{
	name: "null",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("null", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, len(input))

		for k, _ := range args {
			res[k] = ValueOf(reflect.ValueOf(new(any)).Elem())
		}

		return res, nil
	},
}
