package dasel

import (
	"errors"
	"fmt"
)

var ErrMissingPreviousNode = errors.New("missing previous node")

type UnknownComparisonOperatorErr struct {
	Operator string
}

func (e UnknownComparisonOperatorErr) Error() string {
	return fmt.Sprintf("unknown comparison operator: %s", e.Operator)
}

type InvalidIndexErr struct {
	Index string
}

func (e InvalidIndexErr) Error() string {
	return fmt.Sprintf("invalid index: %s", e.Index)
}

type UnsupportedSelector struct {
	Selector string
}

func (e UnsupportedSelector) Error() string {
	return fmt.Sprintf("selector is not supported here: %s", e.Selector)
}

type UnsupportedTypeForSelector struct {
	Selector Selector
	Value    interface{}
}

func (e UnsupportedTypeForSelector) Error() string {
	return fmt.Sprintf("selector [%s] does not support value: %T: %v", e.Selector.Type, e.Value, e.Value)
}

type NotFound struct {
	Selector string
	Node     *Node
}

func (e NotFound) Error() string {
	return fmt.Sprintf("nothing found for selector: %s", e.Selector)
}

type UnexpectedPreviousNilValue struct {
	Selector string
}

func (e UnexpectedPreviousNilValue) Error() string {
	return fmt.Sprintf("previous value is nil: %s", e.Selector)
}
