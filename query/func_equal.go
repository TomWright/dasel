package query

import (
	"fmt"
)

var EqualFunc = BasicFunction{
	name: "equal",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := c.inputValue(s)

		type comparison struct {
			selector string
			value    string
		}

		comparisons := make([]comparison, 0)

		currentComparison := comparison{}

		for i, v := range args {
			switch i % 2 {
			case 0:
				currentComparison.selector = v
			case 1:
				currentComparison.value = v
				comparisons = append(comparisons, currentComparison)
				currentComparison = comparison{}
			}
		}

		match := func(target Value, cmp comparison) (Values, error) {
			matchC := c.subContext(target, cmp.selector)
			finalStep, err := matchC.Run()
			if err != nil {
				return nil, err
			}
			res := make(Values, 0)
			for _, o := range finalStep.output {
				stringExp := cmp.value
				stringGot := fmt.Sprint(o.Interface())
				if stringExp == stringGot {
					res = append(res, target)
				}
			}
			return res, nil
		}

		res := make(Values, 0)

		for _, val := range input {
			for _, cmp := range comparisons {
				vals, err := match(val, cmp)
				if err != nil {
					return nil, err
				}
				res = append(res, vals...)
			}
		}

		return res, nil
	},
}
