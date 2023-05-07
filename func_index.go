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
	o, ok := other.(*ErrIndexNotFound)
	if !ok {
		return false
	}
	if o.Index >= 0 && o.Index != e.Index {
		return false
	}
	return true
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
				case reflect.String:
					runes := []rune(val.String())
					if index < 0 || index > len(runes)-1 {
						if isOptional {
							continue
						}
						return nil, fmt.Errorf("index out of range: %w", &ErrIndexNotFound{Index: index})
					}
					res = append(res, ValueOf(string(runes[index])))
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
					return nil, fmt.Errorf("cannot use index selector on non slice/array types: %w", &ErrIndexNotFound{Index: index})
				}
			}
		}

		return res, nil
	},
	alternativeSelectorFn: func(part string) *Selector {
		if part != "[]" && strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]") {
			strings.Split(strings.TrimPrefix(strings.TrimSuffix(part, "]"), "["), ",")
			return &Selector{
				funcName: "index",
				funcArgs: strings.Split(strings.TrimPrefix(strings.TrimSuffix(part, "]"), "["), ","),
			}
		}
		return nil
	},
}
