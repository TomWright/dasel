package model

import "fmt"

// GoValue returns the value as a native Go value.
func (v *Value) GoValue() (any, error) {
	var res any
	var err error

	switch v.Type() {
	case TypeString:
		res, err = v.StringValue()
	case TypeInt:
		res, err = v.IntValue()
	case TypeFloat:
		res, err = v.FloatValue()
	case TypeBool:
		res, err = v.BoolValue()
	case TypeMap:
		m := make(map[string]any)
		err = v.RangeMap(func(k string, v *Value) error {
			val, err := v.GoValue()
			if err != nil {
				return err
			}
			m[k] = val
			return nil
		})
		res = m
	case TypeSlice:
		s := make([]any, 0)
		err = v.RangeSlice(func(i int, v *Value) error {
			val, err := v.GoValue()
			if err != nil {
				return err
			}
			s = append(s, val)
			return nil
		})
		res = s
	case TypeUnknown:
		res = nil
		err = fmt.Errorf("cannot convert unknown type to Go value")
	case TypeNull:
		res = nil
	default:
		err = fmt.Errorf("unhandled type %v", v.Type())
	}

	return res, err
}
