package dasel

import (
	"fmt"
)

var LenFunc = BasicFunction{
	name: "len",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		if len(args) > 0 {
			return nil, fmt.Errorf("unexpected last args given")
		}

		res := make(Values, 0)

		for _, val := range input {
			res = append(res, ValueOf(val.Len()))
		}

		return res, nil
	},
}
