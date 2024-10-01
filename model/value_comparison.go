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
