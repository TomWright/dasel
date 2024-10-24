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
	return v.UnpackKinds(reflect.Interface, reflect.Ptr).isSlice()
}

func (v *Value) isSlice() bool {
	return v.Value.Kind() == reflect.Slice
}

// Append appends a value to the slice.
func (v *Value) Append(val *Value) error {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	newVal := reflect.Append(unpacked.Value, val.Value)
	unpacked.Value.Set(newVal)
	return nil
}

// SliceLen returns the length of the slice.
func (v *Value) SliceLen() (int, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return 0, ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	return unpacked.Value.Len(), nil
}

// GetSliceIndex returns the value at the specified index in the slice.
func (v *Value) GetSliceIndex(i int) (*Value, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return nil, ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	if i < 0 || i >= unpacked.Value.Len() {
		return nil, SliceIndexOutOfRange{Index: i}
	}
	res := NewValue(unpacked.Value.Index(i))
	return res, nil
}

// SetSliceIndex sets the value at the specified index in the slice.
func (v *Value) SetSliceIndex(i int, val *Value) error {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return ErrUnexpectedType{
			Expected: TypeSlice,
			Actual:   v.Type(),
		}
	}
	if i < 0 || i >= unpacked.Value.Len() {
		return SliceIndexOutOfRange{Index: i}
	}
	unpacked.Value.Index(i).Set(val.Value)
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
