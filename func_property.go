package dasel

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrPropertyNotFound struct {
	Property string
}

func (e ErrPropertyNotFound) Error() string {
	return fmt.Sprintf("property not found: %s", e.Property)
}

func (e ErrPropertyNotFound) Is(other error) bool {
	o, ok := other.(*ErrPropertyNotFound)
	if !ok {
		return false
	}
	if o.Property != "" && o.Property != e.Property {
		return false
	}
	return true
}

var PropertyFunc = BasicFunction{
	name: "property",
	runFn: func(c *Context, s *Step, args []string) (Values, error) {
		if err := requireXOrMoreArgs("property", args, 1); err != nil {
			return nil, err
		}

		input := s.inputs()

		if c.CreateWhenMissing() {
			input = input.initEmptydencodingMaps()
		}

		res := make(Values, 0)

		for _, val := range input {
			for _, property := range args {
				isOptional := strings.HasSuffix(property, "?")
				if isOptional {
					property = strings.TrimSuffix(property, "?")
				}

				switch val.Kind() {
				case reflect.Map:
					index := val.MapIndex(ValueOf(property))
					if index.IsEmpty() {
						if isOptional {
							continue
						}
						if !c.CreateWhenMissing() {
							return nil, fmt.Errorf("could not access map index: %w", &ErrPropertyNotFound{Property: property})
						}
						index = index.asUninitialised()
					}
					res = append(res, index)
				case reflect.Struct:
					value := val.FieldByName(property)
					if value.IsEmpty() {
						if isOptional {
							continue
						}
						return nil, fmt.Errorf("could not access struct field: %w", &ErrPropertyNotFound{Property: property})
					}
					res = append(res, value)
				default:
					if val.IsDencodingMap() {
						index := val.dencodingMapIndex(ValueOf(property))
						if index.IsEmpty() {
							if isOptional {
								continue
							}
							if !c.CreateWhenMissing() {
								return nil, fmt.Errorf("could not access map index: %w", &ErrPropertyNotFound{Property: property})
							}
							index = index.asUninitialised()
						}
						res = append(res, index)
					} else {
						return nil, fmt.Errorf("cannot use property selector on non map/struct types: %s: %w", val.Kind().String(), &ErrPropertyNotFound{Property: property})
					}
				}
			}
		}

		return res, nil
	},
}
