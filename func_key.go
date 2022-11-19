package dasel

var KeyFunc = BasicFunction{
	name: "key",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("key", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, i := range input {
			p := i.Metadata("key")
			if p == nil {
				continue
			}
			res = append(res, ValueOf(p))
		}

		return res, nil
	},
}
