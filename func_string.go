package dasel

import (
	"fmt"
)

var StringFunc = BasicFunction{
	name: "string",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireExactlyXArgs("string", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, len(input))

		for k, v := range args {
			res[k] = ValueOf(fmt.Sprint(v))
		}

		return res, nil
	},
}
