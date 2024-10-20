package model

import (
	"math"
)

func (v *Value) Add(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a + b), nil
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
		return NewValue(a + b), nil
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
		return NewValue(float64(a) + b), nil
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
		return NewValue(a + float64(b)), nil
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
		return NewValue(a + b), nil
	}
	return nil, &ErrIncompatibleTypes{A: v, B: other}
}

func (v *Value) Subtract(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a - b), nil
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
		return NewValue(a - b), nil
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
		return NewValue(float64(a) - b), nil
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
		return NewValue(a - float64(b)), nil
	}
	return nil, &ErrIncompatibleTypes{A: v, B: other}
}

func (v *Value) Multiply(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a * b), nil
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
		return NewValue(a * b), nil
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
		return NewValue(float64(a) * b), nil
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
		return NewValue(a * float64(b)), nil
	}
	return nil, &ErrIncompatibleTypes{A: v, B: other}
}

func (v *Value) Divide(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a / b), nil
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
		return NewValue(a / b), nil
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
		return NewValue(float64(a) / b), nil
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
		return NewValue(a / float64(b)), nil
	}
	return nil, &ErrIncompatibleTypes{A: v, B: other}
}

func (v *Value) Modulo(other *Value) (*Value, error) {
	if v.IsInt() && other.IsInt() {
		a, err := v.IntValue()
		if err != nil {
			return nil, err
		}
		b, err := other.IntValue()
		if err != nil {
			return nil, err
		}
		return NewValue(a % b), nil
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
		return NewValue(math.Mod(a, b)), nil
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
		return NewValue(math.Mod(float64(a), b)), nil
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
		return NewValue(math.Mod(a, float64(b))), nil
	}
	return nil, &ErrIncompatibleTypes{A: v, B: other}
}
