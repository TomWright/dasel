package dasel

import "reflect"

var TypeFunc = BasicFunction{
	name: "type",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("type", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			resStr := "unknown"

			if val.IsDencodingMap() {
				resStr = "object"
			} else {
				switch val.Kind() {
				case reflect.Slice, reflect.Array:
					resStr = "array"
				case reflect.Map, reflect.Struct:
					resStr = "object"
				case reflect.String:
					resStr = "string"
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
					reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
					reflect.Float32, reflect.Float64:
					resStr = "number"
				case reflect.Bool:
					resStr = "bool"
				}
			}
			res = append(res, ValueOf(resStr))
		}

		return res, nil
	},
}
