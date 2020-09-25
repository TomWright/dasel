package dasel

import (
	"fmt"
	"reflect"
)

// Query uses the given selector to query the current node and return the result.
func (n *Node) Query(selector string) (*Node, error) {
	n.Selector.Remaining = selector
	rootNode := n

	if err := buildFindChain(rootNode); err != nil {
		return nil, err
	}

	return lastNode(rootNode), nil
}

func lastNode(n *Node) *Node {
	for {
		if n.Next == nil {
			return n
		}
		n = n.Next
	}
}

func buildFindChain(n *Node) error {
	if n.Selector.Remaining == "" {
		// We've reached the end
		return nil
	}

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
	nextNode.Value, err = findValue(nextNode)
	if err != nil {
		return fmt.Errorf("could not find value: %w", err)
	}

	return buildFindChain(nextNode)
}

// findValueProperty finds the value for the given node using the property selector
// information.
func findValueProperty(n *Node) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Map {
		for _, key := range value.MapKeys() {
			if fmt.Sprint(key.Interface()) == n.Selector.Property {
				return value.MapIndex(key), nil
			}
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value.Type().Kind()}
}

// findValueIndex finds the value for the given node using the index selector
// information.
func findValueIndex(n *Node) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		valueLen := value.Len()
		if n.Selector.Index >= 0 && n.Selector.Index < valueLen {
			return value.Index(n.Selector.Index), nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// findNextAvailableIndex finds the value for the given node using the index selector
// information.
func findNextAvailableIndex(n *Node) (reflect.Value, error) {
	// Next available index isn't supported unless it's creating the item.
	return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
}

// processFindDynamicItem is used by findValueDynamic.
func processFindDynamicItem(n *Node, object reflect.Value) (bool, error) {
	// Loop through each condition.
	allConditionsMatched := true
	for _, c := range n.Selector.Conditions {
		// If the object doesn't match any checks, return a ValueNotFound.
		found, err := c.Check(object)
		if err != nil {
			return false, err
		}
		if !found {
			allConditionsMatched = false
			break
		}
	}
	if allConditionsMatched {
		return true, nil
	}
	return false, nil
}

// findValueDynamic finds the value for the given node using the dynamic selector
// information.
func findValueDynamic(n *Node) (reflect.Value, error) {
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

// findValue finds the value for the given node.
// The value is essentially pulled from the previous node, using the (already parsed) selector
// information stored on the current node.
func findValue(n *Node) (reflect.Value, error) {
	if n.Previous == nil {
		// previous node is required to get it's value.
		return nilValue(), ErrMissingPreviousNode
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return findValueProperty(n)
	case "INDEX":
		return findValueIndex(n)
	case "NEXT_AVAILABLE_INDEX":
		return findNextAvailableIndex(n)
	case "DYNAMIC":
		return findValueDynamic(n)
	default:
		return nilValue(), &UnsupportedSelector{Selector: n.Selector.Type}
	}
}
