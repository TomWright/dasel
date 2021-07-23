package dasel

import (
	"fmt"
	"reflect"
)

// propagate recursively propagates the given nodes value up to the root node.
func propagate(n *Node) error {
	if n.Previous == nil {
		return nil
	}

	if err := propagateValue(n); err != nil {
		return fmt.Errorf("could not propagate value: %w", err)
	}
	return propagate(n.Previous)
}

// propagateValue sends the value of the current node up to the previous node in the chain.
func propagateValue(n *Node) error {
	if n.Previous == nil {
		return nil
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return propagateValueProperty(n)
	case "INDEX":
		return propagateValueIndex(n)
	case "NEXT_AVAILABLE_INDEX":
		return propagateValueNextAvailableIndex(n)
	default:
		return &UnsupportedSelector{Selector: n.Selector.Type}
	}
}

// propagateValueProperty sends the value of the current node up to the previous node in the chain.
func propagateValueProperty(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Map {
		value.SetMapIndex(reflect.ValueOf(n.Selector.Property), n.Value)
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

// propagateValueIndex sends the value of the current node up to the previous node in the chain.
func propagateValueIndex(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		if n.Selector.Index >= 0 && n.Selector.Index < value.Len() {
			value.Index(n.Selector.Index).Set(n.Value)
			return nil
		}
		n.Previous.setReflectValue(reflect.Append(value, n.Value))
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// propagateValueNextAvailableIndex sends the value of the current node up to the previous node in the chain.
func propagateValueNextAvailableIndex(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		n.Previous.setReflectValue(reflect.Append(value, n.Value))
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// deleteFromParent deletes the given node from it's parent.
func deleteFromParent(n *Node) error {
	if n.Previous == nil {
		return nil
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return deleteFromParentProperty(n)
	case "INDEX":
		return deleteFromParentIndex(n)
	default:
		return &UnsupportedSelector{Selector: n.Selector.Type}
	}
}

// deleteFromParentProperty sends the value of the current node up to the previous node in the chain.
func deleteFromParentProperty(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Map {
		value.SetMapIndex(reflect.ValueOf(n.Selector.Property), reflect.Value{})
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

type deletePlaceholder struct {
}

var deletePlaceholderType = reflect.TypeOf(deletePlaceholder{})

// deleteFromParentIndex sends the value of the current node up to the previous node in the chain.
func deleteFromParentIndex(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		if n.Selector.Index >= 0 && n.Selector.Index < value.Len() {
			// Mark this index for deletion.
			// We can't just rewrite the slice here in-case other selectors also target it.
			value.Index(n.Selector.Index).Set(reflect.ValueOf(deletePlaceholder{}))
		}
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// cleanupSliceDeletions scans through the given reflect and removes any invalid reflect values.
// Does not modify the original value.
// Returns false if no modification was made.
func cleanupSliceDeletions(input reflect.Value) (reflect.Value, bool) {
	value := unwrapValue(input)
	if value.Kind() != reflect.Slice {
		return value, false
	}
	res := reflect.MakeSlice(value.Type(), 0, value.Len())

	invalidCount := 0

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		if !item.IsValid() {
			invalidCount++
			continue
		}
		if unwrapValue(item).Type() == deletePlaceholderType {
			invalidCount++
			continue
		}
		res = reflect.Append(res, item)
	}

	if invalidCount == 0 {
		return value, false
	}

	return res, true
}
