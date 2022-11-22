package dasel

import "reflect"

var MergeFunc = BasicFunction{
	name: "merge",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		input := s.inputs()

		res := make(Values, 0)

		if len(args) == 0 {
			// Merge all inputs into a slice.
			resSlice := reflect.MakeSlice(sliceInterfaceType, len(input), len(input))
			for i, val := range input {
				resSlice.Index(i).Set(val.Value)
			}
			resPointer := reflect.New(resSlice.Type())
			resPointer.Elem().Set(resSlice)

			res = append(res, ValueOf(resPointer))
			return res, nil
		}

		// Merge all inputs into a slice.
		resSlice := reflect.MakeSlice(sliceInterfaceType, 0, 0)
		for _, val := range input {
			for _, a := range args {
				gotValues, err := c.subSelect(val, a)
				if err != nil {
					return nil, err
				}

				for _, gotVal := range gotValues {
					resSlice = reflect.Append(resSlice, gotVal.Value)
				}
			}
		}
		resPointer := reflect.New(resSlice.Type())
		resPointer.Elem().Set(resSlice)

		res = append(res, ValueOf(resPointer))
		return res, nil
	},
}
