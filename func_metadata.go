package dasel

var MetadataFunc = BasicFunction{
	name: "metadata",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("metadata", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			for _, a := range args {
				res = append(res, ValueOf(val.Metadata(a)))
			}
		}

		return res, nil
	},
}
