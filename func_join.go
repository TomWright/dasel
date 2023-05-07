package dasel

import (
	"github.com/tomwright/dasel/v2/util"
	"strings"
)

var JoinFunc = BasicFunction{
	name: "join",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("join", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		getValues := func(value Value, selector string) ([]string, error) {
			gotValues, err := c.subSelect(value, selector)
			if err != nil {
				return []string{}, err
			}

			res := make([]string, len(gotValues))
			for k, v := range gotValues {
				res[k] = util.ToString(v.Interface())
			}
			return res, nil
		}

		res := make(Values, 0)

		separator := args[0]
		args = args[1:]

		// No args - join all input values
		if len(args) == 0 {
			values := make([]string, len(input))
			for k, v := range input {
				values[k] = util.ToString(v.Interface())
			}
			res = append(res, ValueOf(strings.Join(values, separator)))
			return res, nil
		}

		// There are args - use each as a selector and join any resulting values.
		values := make([]string, 0)
		for _, val := range input {
			for _, cmp := range args {
				vals, err := getValues(val, cmp)
				if err != nil {
					return nil, err
				}
				values = append(values, vals...)
			}
		}
		res = append(res, ValueOf(strings.Join(values, separator)))

		return res, nil
	},
}
