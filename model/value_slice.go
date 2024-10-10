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
		return fmt.Errorf("expected slice, got %s", v.Type())
	}
	newVal := reflect.Append(unpacked.Value, val.Value)
	unpacked.Value.Set(newVal)
	return nil
}

// SliceLen returns the length of the slice.
func (v *Value) SliceLen() (int, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return 0, fmt.Errorf("expected slice, got %s", v.Type())
	}
	return unpacked.Value.Len(), nil
}

// GetSliceIndex returns the value at the specified index in the slice.
func (v *Value) GetSliceIndex(i int) (*Value, error) {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return nil, fmt.Errorf("expected slice, got %s", v.Type())
	}
	if i < 0 || i >= unpacked.Value.Len() {
		return nil, fmt.Errorf("index out of range: %d", i)
	}
	res := NewValue(unpacked.Value.Index(i))
	return res, nil
}

// SetSliceIndex sets the value at the specified index in the slice.
func (v *Value) SetSliceIndex(i int, val *Value) error {
	unpacked := v.UnpackKinds(reflect.Interface, reflect.Ptr)
	if !unpacked.isSlice() {
		return fmt.Errorf("expected slice, got %s", v.Type())
	}
	if i < 0 || i >= unpacked.Value.Len() {
		return fmt.Errorf("index out of range: %d", i)
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
// If start is -1, it will be treated as 0. e.g. slice[:end] becomes slice[-1:end].
// If end is -1, it will be treated as the length of the slice. e.g. slice[start:] becomes slice[start:-1].
func (v *Value) SliceIndexRange(start, end int) (*Value, error) {
	var err error
	if start == -1 {
		start = 0
	}
	if end == -1 {
		end, err = v.SliceLen()
		if err != nil {
			return nil, fmt.Errorf("error getting slice length: %w", err)
		}
		end = end - 1
		if end < 0 {
			end = 0
		}
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
