package dasel

import (
	"fmt"
	"github.com/tomwright/dasel/internal/storage"
	"reflect"
)

// Put finds the node using the given selector and updates it's value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) Put(selector string, newValue interface{}) error {
	if selector != "." {
		n.Selector.Remaining = selector
	}

	if err := buildPutChain(n); err != nil {
		return err
	}

	final := lastNode(n)

	_, isRealValue := newValue.(storage.RealValue)
	if isRealValue {
		final.setRealValue(newValue)
	} else {
		final.setValue(newValue)
	}

	if final.Selector.Type != "ROOT" {
		if err := propagate(final); err != nil {
			return err
		}
	}

	return nil
}

// PutMultiple all applicable nodes for the given selector and updates all of their values to the given value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) PutMultiple(selector string, newValue interface{}) error {
	if selector != "." {
		n.Selector.Remaining = selector
	}

	if err := buildPutMultipleChain(n); err != nil {
		return err
	}

	final := lastNodes(n)

	val := reflect.ValueOf(newValue)
	_, isRealValue := newValue.(storage.RealValue)

	for _, n := range final {
		if isRealValue {
			n.setRealReflectValue(val)
		} else {
			n.setReflectValue(val)
		}
		if err := propagate(n); err != nil {
			return err
		}
	}

	return nil
}

func buildPutChain(n *Node) error {
	if isFinalSelector(n.Selector.Remaining) {
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
		return fmt.Errorf("could not find put value: %w", err)
	}

	return buildPutChain(nextNode)
}

func buildPutMultipleChain(n *Node) error {
	if isFinalSelector(n.Selector.Remaining) {
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
	n.NextMultiple, err = findNodes(nextSelector, n, true)

	if err != nil {
		return fmt.Errorf("could not find put multiple value: %w", err)
	}

	for _, next := range n.NextMultiple {
		// Add the back reference
		if next.Previous == nil {
			// This can already be set in some cases - SEARCH.
			next.Previous = n
		}

		if err := buildPutMultipleChain(next); err != nil {
			return err
		}
	}

	return nil
}
