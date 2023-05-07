package dasel

import "github.com/tomwright/dasel/v2/util"

var StringFunc = BasicFunction{
	name: "string",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireExactlyXArgs("string", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, len(input))

		for k, v := range args {
			res[k] = ValueOf(util.ToString(v))
		}

		return res, nil
	},
}
