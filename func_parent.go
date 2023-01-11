package dasel

import (
	"strconv"
)

var ParentFunc = BasicFunction{
	name: "parent",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrLessArgs("parent", args, 1); err != nil {
			return nil, err
		}

		levels := 1
		if len(args) > 0 {
			arg, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, err
			}
			levels = arg
		}
		if levels < 1 {
			levels = 1
		}

		input := s.inputs()

		res := make(Values, 0)

		getParent := func(v Value, levels int) (Value, bool) {
			res := v
			for i := 0; i < levels; i++ {
				p := res.Metadata("parent")
				if p == nil {
					return res, false
				}
				if pv, ok := p.(Value); ok {
					res = pv
				} else {
					return res, false
				}
			}
			return res, true
		}

		for _, i := range input {
			if pv, ok := getParent(i, levels); ok {
				res = append(res, pv)
			}
		}

		return res, nil
	},
}
