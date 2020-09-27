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

	if !n.propagated || n.wasInitialised {
		if err := propagateValue(n); err != nil {
			return fmt.Errorf("could not propagate value: %w", err)
		}
	}
	return propagate(n.Previous)
}

// propagateValue sends the value of the current node up the chain.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValue(n *Node) error {
	if n.Previous == nil {
		return nil
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return propagateValueProperty(n)
	case "INDEX":
		return propagateValueIndex(n)
	case "DYNAMIC":
		return propagateValueDynamic(n)
	case "NEXT_AVAILABLE_INDEX":
		return propagateValueNextAvailableIndex(n)
	default:
		return &UnsupportedSelector{Selector: n.Selector.Type}
	}
}

// propagateValueProperty does the opposite of findValueProperty.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueProperty(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Map {
		value.SetMapIndex(reflect.ValueOf(n.Selector.Property), n.Value)
		n.Previous.propagated = true
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

// propagateValueIndex does the opposite of findValueIndex.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
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
		n.Previous.Value = reflect.Append(value, n.Value)
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// propagateValueNextAvailableIndex does the opposite of findValueNextAvailableIndex.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueNextAvailableIndex(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		n.Previous.Value = reflect.Append(value, n.Value)
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// propagateValueDynamic finds the value for the given node using the dynamic selector
// information.
func propagateValueDynamic(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			found, err := processFindDynamicItem(n, object)
			if err != nil {
				return err
			}
			if found {
				object.Set(n.Value)
				return nil
			}
		}
		return &ValueNotFound{Selector: n.Selector.Current, Node: n}
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}
