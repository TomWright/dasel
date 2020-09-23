package dasel

import (
	"errors"
	"fmt"
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

// InvalidIndexErr is returned when a selector targets an index that does not exist.
type InvalidIndexErr struct {
	Index string
}

// Error returns the error message.
func (e InvalidIndexErr) Error() string {
	return fmt.Sprintf("invalid index: %s", e.Index)
}

// UnsupportedSelector is returned when a specific selector type is used in the wrong context.
type UnsupportedSelector struct {
	Selector string
}

// Error returns the error message.
func (e UnsupportedSelector) Error() string {
	return fmt.Sprintf("selector is not supported here: %s", e.Selector)
}

// UnsupportedTypeForSelector is returned when a selector attempts to handle a data type it can't handle.
type UnsupportedTypeForSelector struct {
	Selector Selector
	Value    interface{}
}

// Error returns the error message.
func (e UnsupportedTypeForSelector) Error() string {
	return fmt.Sprintf("selector [%s] does not support value: %T: %v", e.Selector.Type, e.Value, e.Value)
}

// ValueNotFound is returned when a selector string cannot be fully resolved.
type ValueNotFound struct {
	Selector string
	Node     *Node
}

// Error returns the error message.
func (e ValueNotFound) Error() string {
	var previousValue interface{}
	if e.Node != nil && e.Node.Previous != nil {
		previousValue = e.Node.Previous.Value
	}
	return fmt.Sprintf("no value found for selector: %s: %v", e.Selector, previousValue)
}

// UnexpectedPreviousNilValue is returned when the previous node contains a nil value.
type UnexpectedPreviousNilValue struct {
	Selector string
}

// Error returns the error message.
func (e UnexpectedPreviousNilValue) Error() string {
	return fmt.Sprintf("previous value is nil: %s", e.Selector)
}

// UnhandledCheckType is returned when the a check doesn't know how to deal with the given type
type UnhandledCheckType struct {
	Value interface{}
}

// Error returns the error message.
func (e UnhandledCheckType) Error() string {
	return fmt.Sprintf("unhandled check type: %T", e.Value)
}
