package dasel

var ThisFunc = BasicFunction{
	name: "this",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		return c.inputValue(s), nil
	},
}
