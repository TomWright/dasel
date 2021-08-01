package dasel

import "reflect"

// Delete uses the given selector to find and delete the final node from the current node.
func (n *Node) Delete(selector string) error {
	if isFinalSelector(selector) {
		n.setReflectValue(initialiseEmptyOfType(n.Value))
		return nil
	}

	n.Selector.Remaining = selector
	rootNode := n

	if err := buildFindChain(rootNode); err != nil {
		return err
	}

	finalNode := lastNode(rootNode)

	if err := deleteFromParent(finalNode); err != nil {
		return err
	}
	if finalNode.Previous != nil {
		if newSlice, changed := cleanupSliceDeletions(finalNode.Previous.Value); changed {
			finalNode.Previous.setReflectValue(newSlice)
		}

		if finalNode.Previous.Selector.Type != "ROOT" {
			if err := propagate(finalNode.Previous); err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteMultiple uses the given selector to query the current node for every match
// possible and deletes them from the current node.
func (n *Node) DeleteMultiple(selector string) error {
	if selector == "." {
		n.setReflectValue(initialiseEmptyOfType(n.Value))
		return nil
	}

	n.Selector.Remaining = selector

	if err := buildFindMultipleChain(n); err != nil {
		return err
	}

	lastNodes := lastNodes(n)
	for _, lastNode := range lastNodes {
		// delete properties and mark indexes for deletion
		if err := deleteFromParent(lastNode); err != nil {
			return err
		}
	}

	for _, lastNode := range lastNodes {
		// Cleanup indexes marked for deletion
		if lastNode.Previous != nil {
			if newSlice, changed := cleanupSliceDeletions(lastNode.Previous.Value); changed {
				lastNode.Previous.setReflectValue(newSlice)
			}
		}
	}
	for _, lastNode := range lastNodes {
		// Propagate values
		if lastNode.Previous != nil {
			if lastNode.Previous.Selector.Type != "ROOT" {
				if err := propagate(lastNode.Previous); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func initialiseEmptyOfType(value reflect.Value) reflect.Value {
	value = unwrapValue(value)
	switch value.Kind() {
	case reflect.Slice:
		return reflect.MakeSlice(value.Type(), 0, 0)
	case reflect.Map:
		return reflect.MakeMap(value.Type())
	default:
		return reflect.ValueOf(map[string]interface{}{})
	}
}
