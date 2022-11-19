package dasel

import (
	"fmt"
)

var FilterFunc = BasicFunction{
	name: "filter",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("filter", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		runComparison := func(value Value, selector string) (bool, error) {
			gotValues, err := c.subSelect(value, selector)
			if err != nil {
				return false, err
			}

			if len(gotValues) > 1 {
				return false, fmt.Errorf("filter expects selector to return a single value")
			}

			if len(gotValues) == 0 {
				return false, nil
			}

			return IsTruthy(gotValues[0]), nil
		}

		res := make(Values, 0)

		for _, val := range input {
			valPassed := true
			for _, cmp := range args {
				pass, err := runComparison(val, cmp)
				if err != nil {
					return nil, err
				}
				if !pass {
					valPassed = false
					break
				}
			}
			if valPassed {
				res = append(res, val)
			}
		}

		return res, nil
	},
}
