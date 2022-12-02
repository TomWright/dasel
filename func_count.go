package dasel

var CountFunc = BasicFunction{
	name: "count",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := s.inputs()

		return Values{
			ValueOf(len(input)),
		}, nil
	},
}
