package query

import (
	"fmt"
)

var MetadataFunc = BasicFunction{
	name: "metadata",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		if len(args) == 0 {
			return nil, fmt.Errorf("unexpected metadata args given")
		}

		res := make(Values, 0)

		for _, val := range input {
			for _, a := range args {
				res = append(res, ValueOf(val.Metadata(a)))
			}
		}

		return res, nil
	},
}
