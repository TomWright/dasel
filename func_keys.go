package dasel

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type ErrInvalidType struct {
	ExpectedTypes []string
	CurrentType   string
}

func (e *ErrInvalidType) Error() string {
	return fmt.Sprintf("unexpected types: expect %s, get %s", strings.Join(e.ExpectedTypes, " "), e.CurrentType)
}

func (e *ErrInvalidType) Is(other error) bool {
	o, ok := other.(*ErrInvalidType)
	if !ok {
		return false
	}
	if len(e.ExpectedTypes) != len(o.ExpectedTypes) {
		return false
	}
	if e.CurrentType != o.CurrentType {
		return false
	}
	for i, t := range e.ExpectedTypes {
		if t != o.ExpectedTypes[i] {
			return false
		}
	}
	return true
}

var KeysFunc = BasicFunction{
	name: "keys",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireNoArgs("keys", args); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, len(input))

		for i, val := range input {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				list := make([]any, 0, val.Len())

				for i := 0; i < val.Len(); i++ {
					list = append(list, i)
				}

				res[i] = ValueOf(list)
			case reflect.Map:
				keys := val.MapKeys()

				// we expect map keys to be string first so that we can sort them
				list, ok := getStringList(keys)
				if !ok {
					list = getAnyList(keys)
				}

				res[i] = ValueOf(list)
			default:
				return nil, &ErrInvalidType{
					ExpectedTypes: []string{"slice", "array", "map"},
					CurrentType:   val.Kind().String(),
				}
			}
		}

		return res, nil
	},
}

func getStringList(values []Value) ([]any, bool) {
	stringList := make([]string, len(values))
	for i, v := range values {
		if v.Kind() != reflect.String {
			return nil, false
		}
		stringList[i] = v.String()
	}

	sort.Strings(stringList)

	anyList := make([]any, len(stringList))
	for i, v := range stringList {
		anyList[i] = v
	}

	return anyList, true
}

func getAnyList(values []Value) []any {
	anyList := make([]any, len(values))
	for i, v := range values {
		anyList[i] = v.Interface()
	}
	return anyList
}
