package dasel

import (
	"fmt"
	"reflect"
)

var NotFunc = BasicFunction{
	name: "not",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("not", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		runComparison := func(value Value, selector string) (bool, error) {
			gotValues, err := c.subSelect(value, selector)
			if err != nil {
				return false, err
			}

			if len(gotValues) > 1 {
				return false, fmt.Errorf("not expects selector to return a single value")
			}

			if len(gotValues) == 0 {
				return false, nil
			}

			return IsTruthy(gotValues[0].Interface()), nil
		}

		res := make(Values, 0)

		for _, val := range input {
			for _, selector := range args {
				truthy, err := runComparison(val, selector)
				if err != nil {
					return nil, err
				}
				res = append(res, Value{Value: reflect.ValueOf(!truthy)})
			}
		}

		return res, nil
	},
}
