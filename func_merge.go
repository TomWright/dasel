package dasel

var MergeFunc = BasicFunction{
	name: "merge",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("merge", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		type comparison struct {
			selector string
			value    string
		}

		res := make(Values, 0)

		for _, val := range input {
			for _, a := range args {
				gotValues, err := c.subSelect(val, a)
				if err != nil {
					return nil, err
				}

				res = append(res, gotValues...)
			}
		}

		return res, nil
	},
}
