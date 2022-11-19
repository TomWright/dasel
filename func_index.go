package dasel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ErrIndexNotFound struct {
	Index int
}

func (e ErrIndexNotFound) Error() string {
	return fmt.Sprintf("index not found: %d", e.Index)
}

func (e ErrIndexNotFound) Is(other error) bool {
	_, ok := other.(*ErrIndexNotFound)
	return ok
}

var IndexFunc = BasicFunction{
	name: "index",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("index", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		res := make(Values, 0)

		for _, val := range input {
			for _, indexStr := range args {
				isOptional := strings.HasSuffix(indexStr, "?")
				if isOptional {
					indexStr = strings.TrimSuffix(indexStr, "?")
				}

				index, err := strconv.Atoi(indexStr)
				if err != nil {
					if isOptional {
						continue
					}
					return nil, fmt.Errorf("invalid index: %w", err)
				}

				switch val.Kind() {
				case reflect.Slice, reflect.Array:
					if index < 0 || index > val.Len()-1 {
						if isOptional {
							continue
						}
						return nil, fmt.Errorf("index out of range: %w", &ErrIndexNotFound{Index: index})
					}
					value := val.Index(index)
					res = append(res, value)
				default:
					return nil, fmt.Errorf("cannot use index selector on non slice/array types")
				}
			}
		}

		return res, nil
	},
	alternativeSelectorFn: func(part string) *Selector {
		if part != "[]" && strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			return &Selector{
				funcName: "index",
				funcArgs: []string{
					strings.TrimPrefix(strings.TrimSuffix(part, "]"), "["),
				},
			}
		}
		return nil
	},
}
