package dasel

var ThisFunc = BasicFunction{
	name: "this",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("this", args); err != nil {
			return nil, err
		}
		return s.inputs(), nil
	},
}
