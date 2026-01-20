package model

import (
	"fmt"
	"reflect"
)

// NewSliceValue returns a new slice value.
func NewSliceValue() *Value {
	res := newPtr()
	s := reflect.MakeSlice(reflect.SliceOf(reflect.TypeFor[any]()), 0, 0)
	ptr := reflect.New(reflect.SliceOf(reflect.TypeFor[any]()))
	ptr.Elem().Set(s)
	res.Elem().Set(ptr)
	return NewValue(res)
}

// IsSlice returns true if the value is a slice.
func (v *Value) IsSlice() bool {
	return v.UnpackKinds(reflect.Interface, reflect.Pointer).isSlice()
}

func (v *Value) isSlice() bool {
	return v.value.Kind() == reflect.Slice
}

// Append appends a value to the slice.
func (v *Value) Append(val *Value) error {
	// Branches behave differently when appending to a slice.
	// We expect each item in a branch to be its own value.
	if val.IsBranch() {
		return val.RangeSlice(func(_ int, item *Value) error {
			return v.Append(item)
		})
	}

	unpacked := v.UnpackKinds(reflect.Interface, reflect.Pointer)
	if !unpacked.isSlice() {
		return ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}

	valToAppend := val.value
	// Wrap the value if it has a set function. This helps ensures sets are made correctly.
	// This was first noticed as an issue with the recursive descent response slice.
	// We could always wrap, but that would be less efficient.
	if val.isDaselValue() || val.setFn != nil || val.IsScalar() {
		valToAppend = reflect.ValueOf(val)
	}
	newVal := reflect.Append(unpacked.value, valToAppend)
	unpacked.value.Set(newVal)
	return nil
}

// SliceLen returns the length of the slice.
func (v *Value) SliceLen() (int, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Pointer)
	if !unpacked.isSlice() {
		return 0, ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	return unpacked.value.Len(), nil
}

// GetSliceIndex returns the value at the specified index in the slice.
func (v *Value) GetSliceIndex(i int) (*Value, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Pointer)
	if !unpacked.isSlice() {
		return nil, ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	if i < 0 || i >= unpacked.value.Len() {
		return nil, SliceIndexOutOfRange{Index: i}
	}

	item := unpacked.value.Index(i)
	if item.Kind() == reflect.Pointer && item.Type() == reflect.TypeFor[*Value]() {
		return item.Interface().(*Value), nil
	}

	res := NewValue(item)
	return res, nil
}

// SetSliceIndex sets the value at the specified index in the slice.
func (v *Value) SetSliceIndex(i int, val *Value) error {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Pointer)
	if !unpacked.isSlice() {
		return ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	if i < 0 || i >= unpacked.value.Len() {
		return SliceIndexOutOfRange{Index: i}
	}
	unpacked.value.Index(i).Set(reflect.ValueOf(val))
	return nil
}

// RangeSlice iterates over each item in the slice and calls the provided function.
func (v *Value) RangeSlice(f func(int, *Value) error) error {
	length, err := v.SliceLen()
	if err != nil {
		return fmt.Errorf("error getting slice length: %w", err)
	}

	for i := 0; i < length; i++ {
		va, err := v.GetSliceIndex(i)
		if err != nil {
			return fmt.Errorf("error getting slice index %d: %w", i, err)
		}
		if err := f(i, va); err != nil {
			return err
		}
	}

	return nil
}

// SliceIndexRange returns a new slice containing the values between the start and end indexes.
// Comparable to go's slice[start:end].
func (v *Value) SliceIndexRange(start, end int) (*Value, error) {
	l, err := v.SliceLen()
	if err != nil {
		return nil, fmt.Errorf("error getting slice length: %w", err)
	}

	if start < 0 {
		start = l + start
	}
	if end < 0 {
		end = l + end
	}

	res := NewSliceValue()

	if start > end {
		for i := start; i >= end; i-- {
			item, err := v.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			if err := res.Append(item); err != nil {
				return nil, fmt.Errorf("error appending value to slice: %w", err)
			}
		}
	} else {
		for i := start; i <= end; i++ {
			item, err := v.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			if err := res.Append(item); err != nil {
				return nil, fmt.Errorf("error appending value to slice: %w", err)
			}
		}
	}

	return res, nil
}
