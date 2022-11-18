package dasel

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrMissingPreviousNode is returned when findValue doesn't have access to the previous node.
var ErrMissingPreviousNode = errors.New("missing previous node")

// UnknownComparisonOperatorErr is returned when
type UnknownComparisonOperatorErr struct {
	Operator string
}

// Error returns the error message.
func (e UnknownComparisonOperatorErr) Error() string {
	return fmt.Sprintf("unknown comparison operator: %s", e.Operator)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e UnknownComparisonOperatorErr) Is(err error) bool {
	_, ok := err.(*UnknownComparisonOperatorErr)
	return ok
}

// InvalidIndexErr is returned when a selector targets an index that does not exist.
type InvalidIndexErr struct {
	Index string
}

// Error returns the error message.
func (e InvalidIndexErr) Error() string {
	return fmt.Sprintf("invalid index: %s", e.Index)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e InvalidIndexErr) Is(err error) bool {
	_, ok := err.(*InvalidIndexErr)
	return ok
}

// UnsupportedSelector is returned when a specific selector type is used in the wrong context.
type UnsupportedSelector struct {
	Selector string
}

// Error returns the error message.
func (e UnsupportedSelector) Error() string {
	return fmt.Sprintf("selector is not supported here: %s", e.Selector)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e UnsupportedSelector) Is(err error) bool {
	_, ok := err.(*UnsupportedSelector)
	return ok
}

// ValueNotFound is returned when a selector string cannot be fully resolved.
type ValueNotFound struct {
	Selector      string
	PreviousValue reflect.Value
}

// Error returns the error message.
func (e ValueNotFound) Error() string {
	return fmt.Sprintf("no value found for selector: %s: %v", e.Selector, e.PreviousValue)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e ValueNotFound) Is(err error) bool {
	_, ok := err.(*ValueNotFound)
	return ok
}

// UnexpectedPreviousNilValue is returned when the previous node contains a nil value.
type UnexpectedPreviousNilValue struct {
	Selector string
}

// Error returns the error message.
func (e UnexpectedPreviousNilValue) Error() string {
	return fmt.Sprintf("previous value is nil: %s", e.Selector)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e UnexpectedPreviousNilValue) Is(err error) bool {
	_, ok := err.(*UnexpectedPreviousNilValue)
	return ok
}

// UnhandledCheckType is returned when the a check doesn't know how to deal with the given type
type UnhandledCheckType struct {
	Value interface{}
}

// Error returns the error message.
func (e UnhandledCheckType) Error() string {
	return fmt.Sprintf("unhandled check type: %T", e.Value)
}

// Is implements the errors interface, so the errors.Is() function can be used.
func (e UnhandledCheckType) Is(err error) bool {
	_, ok := err.(*UnhandledCheckType)
	return ok
}
