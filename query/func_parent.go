package query

var ParentFunc = BasicFunction{
	name: "parent",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		res := make(Values, 0)

		for _, i := range input {
			p := i.Metadata("parent")
			if p == nil {
				continue
			}
			if pv, ok := p.(Value); ok {
				res = append(res, pv)
			}
		}

		return res, nil
	},
}
