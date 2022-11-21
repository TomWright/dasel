package dasel

var MergeFunc = BasicFunction{
	name: "merge",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := s.inputs()

		type comparison struct {
			selector string
			value    string
		}

		res := make(Values, 0)

		if len(args) == 0 {
			res = append(res, input...)
			return res, nil
		}

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
