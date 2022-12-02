package dasel

import (
	"fmt"
	"reflect"
)

var MapOfFunc = BasicFunction{
	name: "mapOf",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("mapOf", args, 2); err != nil {
			return nil, err
		}
		if err := requireModulusXArgs("mapOf", args, 2); err != nil {
			return nil, err
		}

		input := s.inputs()

		type pair struct {
			key      string
			selector string
		}

		pairs := make([]pair, 0)

		currentPair := pair{}

		for i, v := range args {
			switch i % 2 {
			case 0:
				currentPair.key = v
			case 1:
				currentPair.selector = v
				pairs = append(pairs, currentPair)
				currentPair = pair{}
			}
		}

		getValue := func(value Value, p pair) (Value, error) {
			gotValues, err := c.subSelect(value, p.selector)
			if err != nil {
				return Value{}, err
			}

			if len(gotValues) != 1 {
				return Value{}, fmt.Errorf("mapOf expects selector to return exactly 1 value")
			}

			return gotValues[0], nil
		}

		res := make(Values, 0)

		for _, val := range input {
			result := reflect.MakeMap(mapStringInterfaceType)

			for _, p := range pairs {
				gotValue, err := getValue(val, p)
				if err != nil {
					return nil, err
				}

				result.SetMapIndex(reflect.ValueOf(p.key), gotValue.Value)
			}

			res = append(res, ValueOf(result))
		}

		return res, nil
	},
}
