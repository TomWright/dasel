package model

func (v *Value) Equal(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a == b), nil
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
		return NewValue(a == b), nil
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
	return nil, &ErrIncompatibleTypes{A: v, B: other}
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

func (v *Value) EqualTypeValue(other *Value) bool {
	if v.Type() != other.Type() {
		return false
	}

	switch v.Type() {
	case TypeString:
		a, _ := v.StringValue()
		b, _ := other.StringValue()
		return a == b
	case TypeInt:
		a, _ := v.IntValue()
		b, _ := other.IntValue()
		return a == b
	case TypeFloat:
		a, _ := v.FloatValue()
		b, _ := other.FloatValue()
		return a == b
	case TypeBool:
		a, _ := v.BoolValue()
		b, _ := other.BoolValue()
		return a == b
	case TypeMap:
		a, _ := v.MapKeys()
		b, _ := other.MapKeys()
		if len(a) != len(b) {
			return false
		}
		for _, key := range a {
			valA, _ := v.GetMapKey(key)
			valB, _ := other.GetMapKey(key)
			if !valA.EqualTypeValue(valB) {
				return false
			}
		}
		return true
	case TypeSlice:
		a, _ := v.SliceLen()
		b, _ := other.SliceLen()
		if a != b {
			return false
		}
		for i := 0; i < a; i++ {
			valA, _ := v.GetSliceIndex(i)
			valB, _ := other.GetSliceIndex(i)
			if !valA.EqualTypeValue(valB) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
