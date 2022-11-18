package dasel

import (
	"fmt"
	"reflect"
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

		runComparison := func(value Value, cmp comparison) (bool, error) {
			gotValues, err := performSubQuery(c, value, cmp.selector)
			if err != nil {
				return false, err
			}

			if len(gotValues) > 1 {
				return false, fmt.Errorf("equal expects selector to return a single value")
			}

			if len(gotValues) == 0 {
				return false, nil
			}

			gotValue := fmt.Sprint(gotValues[0].Interface())
			return gotValue == cmp.value, nil
		}

		res := make(Values, 0)

		for _, val := range input {
			valPassed := true
			for _, cmp := range comparisons {
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
