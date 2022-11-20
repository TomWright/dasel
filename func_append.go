package dasel

import (
	"fmt"
	"reflect"
)

var AppendFunc = BasicFunction{
	name: "append",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("append", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		if c.CreateWhenMissing() {
			input = input.initEmptySlices()
		}

		res := make(Values, 0)

		for _, val := range input {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				val = val.Append()
				value := val.Index(val.Len() - 1)
				res = append(res, value)
			default:
				return nil, fmt.Errorf("cannot use append selector on non slice/array types")
			}
		}

		return res, nil
	},
	alternativeSelectorFn: func(part string) *Selector {
		if part == "[]" {
			return &Selector{
				funcName: "append",
				funcArgs: []string{},
			}
		}
		return nil
	},
}
