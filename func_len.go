package dasel

var LenFunc = BasicFunction{
	name: "len",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("len", args); err != nil {
			return nil, err
		}

		input := c.inputValue(s)

		res := make(Values, 0)

		for _, val := range input {
			res = append(res, ValueOf(val.Len()))
		}

		return res, nil
	},
}
