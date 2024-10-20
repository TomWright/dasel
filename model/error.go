package model

import "fmt"

// MapKeyNotFound is returned when a key is not found in a map.
type MapKeyNotFound struct {
	Key string
}

// Error returns the error message.
func (e *MapKeyNotFound) Error() string {
	return fmt.Sprintf("map key not found: %q", e.Key)
}

// SliceIndexOutOfRange is returned when an index is invalid.
type SliceIndexOutOfRange struct {
	Index int
}

// Error returns the error message.
func (e *SliceIndexOutOfRange) Error() string {
	return fmt.Sprintf("slice index out of range: %d", e.Index)
}

// ErrIncompatibleTypes is returned when two values are incompatible.
type ErrIncompatibleTypes struct {
	A *Value
	B *Value
}

// Error returns the error message.
func (e *ErrIncompatibleTypes) Error() string {
	return fmt.Sprintf("incompatible types: %s and %s", e.A.Type(), e.B.Type())
}
