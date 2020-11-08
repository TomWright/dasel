package dasel

import (
	"fmt"
	"reflect"
)

// Put finds the node using the given selector and updates it's value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) Put(selector string, newValue interface{}) error {
	n.Selector.Remaining = selector

	if err := buildPutChain(n); err != nil {
		return err
	}

	final := lastNode(n)

	final.Value = reflect.ValueOf(newValue)

	if err := propagate(final); err != nil {
		return err
	}

	return nil
}

// PutMultiple all applicable nodes for the given selector and updates all of their values to the given value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) PutMultiple(selector string, newValue interface{}) error {
	n.Selector.Remaining = selector

	if err := buildPutMultipleChain(n); err != nil {
		return err
	}

	final := lastNodes(n)

	val := reflect.ValueOf(newValue)

	for _, n := range final {
		n.Value = val
		if err := propagate(n); err != nil {
			return err
		}
	}

	return nil
}

func buildPutChain(n *Node) error {
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
	n.Next = []*Node{nextNode}
	nextNode.Previous = n

	// Populate the value for the new node.
	nextNode.Value, err = findValue(nextNode, true)
	if err != nil {
		return fmt.Errorf("could not find put value: %w", err)
	}

	return buildPutChain(nextNode)
}

func buildPutMultipleChain(n *Node) error {
	if n.Selector.Remaining == "" {
		// We've reached the end
		return nil
	}

	var err error

	// Parse the selector.
	nextSelector, err := ParseSelector(n.Selector.Remaining)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %w", err)
	}

	// Populate the value for the new node.
	n.Next, err = findNodes(nextSelector, n.Value, true)

	if err != nil {
		return fmt.Errorf("could not find put multiple value: %w", err)
	}

	for _, next := range n.Next {
		// Add the back reference
		next.Previous = n

		if err := buildPutMultipleChain(next); err != nil {
			return err
		}
	}

	return nil
}
