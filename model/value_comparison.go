package model

func (v *Value) Compare(other *Value) (int, error) {
	eq, err := v.Equal(other)
	if err != nil {
		return 0, err
	}
	eqVal, err := eq.BoolValue()
	if err != nil {
		return 0, err
	}
	if eqVal {
		return 0, nil
	}

	lt, err := v.LessThan(other)
	if err != nil {
		return 0, err
	}
	ltVal, err := lt.BoolValue()
	if err != nil {
		return 0, err
	}
	if ltVal {
		return -1, nil
	}

	return 1, nil
}

func (v *Value) Equal(other *Value) (*Value, error) {
	if v.IsInt() && other.IsFloat() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.FloatValue()
		if err != nil {
			return nil, err
		}
		return NewValue(float64(a) == b), nil
	}
	if v.IsFloat() && other.IsInt() {
		a, err := v.FloatValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a == float64(b)), nil
	}

	if v.Type() != other.Type() {
		return nil, &ErrIncompatibleTypes{A: v, B: other}
	}

	isEqual, err := v.EqualTypeValue(other)
	if err != nil {
		return nil, err
	}
	return NewValue(isEqual), nil
}

func (v *Value) NotEqual(other *Value) (*Value, error) {
	equals, err := v.Equal(other)
	if err != nil {
		return nil, err
	}
	boolValue, err := equals.BoolValue()
	if err != nil {
		return nil, err
	}
	return NewValue(!boolValue), nil
}

func (v *Value) LessThan(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a < b), nil
	}
	if v.IsFloat() && other.IsFloat() {
		a, err := v.FloatValue()
		if err != nil {
			return nil, err
		}
		b, err := other.FloatValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a < b), nil
	}
	if v.IsInt() && other.IsFloat() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.FloatValue()
		if err != nil {
			return nil, err
		}
		return NewValue(float64(a) < b), nil
	}
	if v.IsFloat() && other.IsInt() {
		a, err := v.FloatValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a < float64(b)), nil
	}

	if v.IsString() && other.IsString() {
		a, err := v.StringValue()
		if err != nil {
			return nil, err
		}
		b, err := other.StringValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a < b), nil
	}

	return nil, &ErrIncompatibleTypes{A: v, B: other}
}

func (v *Value) LessThanOrEqual(other *Value) (*Value, error) {
	lessThan, err := v.LessThan(other)
	if err != nil {
		return nil, err
	}
	boolValue, err := lessThan.BoolValue()
	if err != nil {
		return nil, err
	}
	equals, err := v.Equal(other)
	if err != nil {
		return nil, err
	}
	boolEquals, err := equals.BoolValue()
	if err != nil {
		return nil, err
	}
	return NewValue(boolValue || boolEquals), nil
}

func (v *Value) GreaterThan(other *Value) (*Value, error) {
	lessThanOrEqual, err := v.LessThanOrEqual(other)
	if err != nil {
		return nil, err
	}
	boolValue, err := lessThanOrEqual.BoolValue()
	if err != nil {
		return nil, err
	}
	return NewValue(!boolValue), nil
}

func (v *Value) GreaterThanOrEqual(other *Value) (*Value, error) {
	lessThan, err := v.LessThan(other)
	if err != nil {
		return nil, err
	}
	boolValue, err := lessThan.BoolValue()
	if err != nil {
		return nil, err
	}
	return NewValue(!boolValue), nil
}

func (v *Value) EqualTypeValue(other *Value) (bool, error) {
	if v.Type() != other.Type() {
		return false, nil
	}

	switch v.Type() {
	case TypeString:
		a, err := v.StringValue()
		if err != nil {
			return false, err
		}
		b, err := other.StringValue()
		if err != nil {
			return false, err
		}
		return a == b, nil
	case TypeInt:
		a, err := v.IntValue()
		if err != nil {
			return false, err
		}
		b, err := other.IntValue()
		if err != nil {
			return false, err
		}
		return a == b, nil
	case TypeFloat:
		a, err := v.FloatValue()
		if err != nil {
			return false, err
		}
		b, err := other.FloatValue()
		if err != nil {
			return false, err
		}
		return a == b, nil
	case TypeBool:
		a, err := v.BoolValue()
		if err != nil {
			return false, err
		}
		b, err := other.BoolValue()
		if err != nil {
			return false, err
		}
		return a == b, nil
	case TypeMap:
		a, err := v.MapKeys()
		if err != nil {
			return false, err
		}
		b, err := other.MapKeys()
		if err != nil {
			return false, err
		}
		if len(a) != len(b) {
			return false, nil
		}
		for _, key := range a {
			valA, err := v.GetMapKey(key)
			if err != nil {
				return false, err
			}
			valB, err := other.GetMapKey(key)
			if err != nil {
				return false, err
			}
			equal, err := valA.EqualTypeValue(valB)
			if err != nil {
				return false, err
			}
			if !equal {
				return false, nil
			}
		}
		return true, nil
	case TypeSlice:
		a, err := v.SliceLen()
		if err != nil {
			return false, err
		}
		b, err := other.SliceLen()
		if err != nil {
			return false, err
		}
		if a != b {
			return false, nil
		}
		for i := 0; i < a; i++ {
			valA, err := v.GetSliceIndex(i)
			if err != nil {
				return false, err
			}
			valB, err := other.GetSliceIndex(i)
			if err != nil {
				return false, err
			}
			equal, err := valA.EqualTypeValue(valB)
			if err != nil {
				return false, err
			}
			if !equal {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, nil
	}
}
