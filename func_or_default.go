package dasel

import (
	"errors"
	"fmt"
)

var OrDefaultFunc = BasicFunction{
	name: "orDefault",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireExactlyXArgs("orDefault", args, 2); err != nil {
			return nil, err
		}

		input := s.inputs()

		if c.CreateWhenMissing() {
			input = input.initEmptydencodingMaps()
		}

		runSubselect := func(value Value, selector string, defaultSelector string) (Value, error) {
			gotValues, err := c.subSelect(value, selector)
			notFound := false
			if err != nil {
				if errors.Is(err, &ErrPropertyNotFound{}) {
					notFound = true
				} else if errors.Is(err, &ErrIndexNotFound{Index: -1}) {
					notFound = true
				} else {
					return Value{}, err
				}
			}

			if !notFound {
				// Check result of first query
				if len(gotValues) != 1 {
					return Value{}, fmt.Errorf("orDefault expects selector to return exactly 1 value")
				}

				// Consider nil values as not found
				if gotValues[0].IsNil() {
					notFound = true
				}
			}

			if notFound {
				gotValues, err = c.subSelect(value, defaultSelector)
				if err != nil {
					return Value{}, err
				}
				if len(gotValues) != 1 {
					return Value{}, fmt.Errorf("orDefault expects selector to return exactly 1 value")
				}
			}

			return gotValues[0], nil
		}

		res := make(Values, 0)

		for _, val := range input {
			resolvedValue, err := runSubselect(val, args[0], args[1])
			if err != nil {
				return nil, err
			}

			res = append(res, resolvedValue)
		}

		return res, nil
	},
}
