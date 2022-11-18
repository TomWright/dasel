package dasel

import (
	"fmt"
	"reflect"
)

var AndFunc = BasicFunction{
	name: "and",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("filter", args, 1); err != nil {
			return nil, err
		}

		input := c.inputValue(s)

		runComparison := func(value Value, selector string) (bool, error) {
			gotValues, err := performSubQuery(c, value, selector)
			if err != nil {
				return false, err
			}

			if len(gotValues) > 1 {
				return false, fmt.Errorf("and expects selector to return a single value")
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
			res = append(res, Value{Value: reflect.ValueOf(valPassed)})
		}

		return res, nil
	},
}
