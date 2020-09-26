package dasel

import (
	"fmt"
	"reflect"
)

// Initialise sets up the node to ensure it has a value.
func (n *Node) Initialise(selector Selector) {
	if !isValid(n.Value) {
		if n.Next == nil {
			var err error
			var val interface{}
			n.Value, err = putValue(n, reflect.ValueOf(val), true)
			if err != nil {
				panic(err)
			}
			return
		}

		// Set an initial value based off of the selector type.
		switch n.Next.Selector.Type {
		case "ROOT":
			n.Value = reflect.ValueOf(make(map[string]interface{}))
		case "PROPERTY":
			var err error
			n.Value, err = putValue(n, reflect.ValueOf(make(map[string]interface{})), true)
			if err != nil {
				panic(err)
			}
		case "NEXT_AVAILABLE_INDEX", "INDEX", "DYNAMIC":
			var err error
			n.Value, err = putValue(n, reflect.ValueOf(make([]interface{}, 0)), true)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Put finds the node using the given selector and updates it's value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) Put(selector string, newValue interface{}) error {
	n.Selector.Remaining = selector

	if err := buildPutChain(n, unwrapValue(reflect.ValueOf(newValue))); err != nil {
		return err
	}

	return nil
}

func buildPutChain(n *Node, newValue reflect.Value) error {
	if n.IsFinal() {
		// We've reached the end
		return nil
	}

	n.Initialise(n.Selector)

	var err error
	nextNode := &Node{}

	// Parse the selector.
	nextNode.Selector, err = ParseSelector(n.Selector.Remaining)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %w", err)
	}

	// Link the nodes.
	n.Next = nextNode
	nextNode.Previous = n

	// Populate the value for the new node.
	nextNode.Value, err = putValue(nextNode, newValue, nextNode.IsFinal())
	if err != nil {
		return fmt.Errorf("could not find value: %w", err)
	}

	return buildPutChain(nextNode, newValue)
}

func putValue(n *Node, newValue reflect.Value, writeNewValue bool) (reflect.Value, error) {
	if n.Previous == nil {
		// previous node is required to get it's value.
		return nilValue(), ErrMissingPreviousNode
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return putValueProperty(n, newValue, writeNewValue)
	case "INDEX":
		return putValueIndex(n, newValue, writeNewValue)
	case "NEXT_AVAILABLE_INDEX":
		return putNextAvailableIndex(n, newValue, writeNewValue)
	case "DYNAMIC":
		return putValueDynamic(n, newValue, writeNewValue)
	default:
		return nilValue(), &UnsupportedSelector{Selector: n.Selector.Type}
	}
}

// putValueProperty writes the new value to the given node using the selector information.
func putValueProperty(n *Node, newValue reflect.Value, writeNewValue bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Map {
		key := reflect.ValueOf(n.Selector.Property)
		if writeNewValue {
			value.SetMapIndex(key, newValue)
		}
		return value.MapIndex(key), nil
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// putValueIndex writes the new value to the given node using the selector information.
func putValueIndex(n *Node, newValue reflect.Value, writeNewValue bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		valueLen := value.Len()
		if n.Selector.Index < 0 || n.Selector.Index >= valueLen {
			// If the requested index isn't within the range of the slice, let's append to it instead.
			return putNextAvailableIndex(n, newValue, writeNewValue)
		}
		val := value.Index(n.Selector.Index)
		if writeNewValue {
			val.Set(newValue)
		}
		return val, nil
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// putNextAvailableIndex writes the new value to the given node using the selector information.
func putNextAvailableIndex(n *Node, newValue reflect.Value, writeNewValue bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		if value.Len() == 0 {
			n.Previous.Value = reflect.Append(value, newValue)
			return newValue, nil
		}
		newValue := reflect.Indirect(reflect.New(value.Index(0).Type()))
		n.Previous.Value = reflect.Append(value, newValue)
		return newValue, nil
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// putValueDynamic writes the new value to the given node using the selector information.
func putValueDynamic(n *Node, newValue reflect.Value, writeNewValue bool) (reflect.Value, error) {
	if writeNewValue {
		return nilValue(), &UnsupportedSelector{Selector: n.Selector.Current}
	}

	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			found, err := processFindDynamicItem(n, object)
			if err != nil {
				return nilValue(), err
			}
			if found {
				return object, nil
			}
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}
