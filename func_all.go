package dasel

import (
	"fmt"
	"github.com/tomwright/dasel/v2/dencoding"
	"reflect"
)

var AllFunc = BasicFunction{
	name: "all",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("all", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			switch val.Kind() {
			case reflect.String:
				for _, r := range val.String() {
					res = append(res, ValueOf(string(r)))
				}
			case reflect.Slice, reflect.Array:
				for i := 0; i < val.Len(); i++ {
					res = append(res, val.Index(i))
				}
			case reflect.Map:
				for _, key := range val.MapKeys() {
					res = append(res, val.MapIndex(key))
				}
			default:
				if val.IsDencodingMap() {
					for _, k := range val.Interface().(*dencoding.Map).Keys() {
						res = append(res, val.dencodingMapIndex(ValueOf(k)))
					}
				} else {
					return nil, fmt.Errorf("cannot use all selector on non slice/array/map types")
				}
			}
		}

		return res, nil
	},
}
