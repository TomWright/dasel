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
	n.Next = nextNode
	nextNode.Previous = n

	// Populate the value for the new node.
	nextNode.Value, err = findValue(nextNode, true)
	if err != nil {
		return fmt.Errorf("could not find value: %w", err)
	}

	return buildPutChain(nextNode)
}
