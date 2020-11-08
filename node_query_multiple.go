package dasel

import (
	"fmt"
	"reflect"
)

// QueryMultiple uses the given selector to query the current node for every match
// possible and returns all of the end nodes.
func (n *Node) QueryMultiple(selector string) ([]*Node, error) {
	n.Selector.Remaining = selector

	if err := buildFindMultipleChain(n); err != nil {
		return nil, err
	}

	return lastNodes(n), nil
}

// lastNodes returns a list of all of the last nodes.
func lastNodes(n *Node) []*Node {
	if n.Next == nil {
		return []*Node{n}
	}
	if len(n.Next) == 0 {
		return []*Node{}
	}
	var res []*Node
	for _, nextNode := range n.Next {
		res = append(res, lastNodes(nextNode)...)
	}
	return res
}

func buildFindMultipleChain(n *Node) error {
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
	n.Next, err = findNodes(nextSelector, n.Value, false)
	if err != nil {
		return fmt.Errorf("could not find multiple value: %w", err)
	}

	for _, next := range n.Next {
		// Add the back reference
		next.Previous = n

		if err := buildFindMultipleChain(next); err != nil {
			return err
		}
	}

	return nil
}

// findNodesProperty finds the value for the given node using the property selector
// information.
func findNodesProperty(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}

	value := unwrapValue(previousValue)

	if value.Kind() == reflect.Map {
		node := &Node{
			Value:    nilValue(),
			Selector: selector,
		}
		for _, key := range value.MapKeys() {
			if fmt.Sprint(key.Interface()) == selector.Property {
				node.Value = value.MapIndex(key)
				return []*Node{node}, nil
			}
		}
		if createIfNotExists {
			return []*Node{node}, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: previousValue.Type().Kind()}
}

// findNodesIndex finds the value for the given node using the index selector
// information.
func findNodesIndex(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}

	value := unwrapValue(previousValue)

	if value.Kind() == reflect.Slice {
		node := &Node{
			Value:    nilValue(),
			Selector: selector,
		}
		valueLen := value.Len()
		if selector.Index >= 0 && selector.Index < valueLen {
			node.Value = value.Index(selector.Index)
			return []*Node{node}, nil
		}
		if createIfNotExists {
			return []*Node{node}, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value.Kind()}
}

// findNextAvailableNodes finds the value for the given node using the index selector
// information.
func findNextAvailableNodes(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !createIfNotExists {
		// Next available index isn't supported unless it's creating the item.
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}
	}
	return []*Node{
		{
			Value:    nilValue(),
			Selector: selector,
		},
	}, nil
}

// processFindDynamicItems is used by findNodesDynamic.
func processFindDynamicItems(selector Selector, object reflect.Value) (bool, error) {
	// Loop through each condition.
	allConditionsMatched := true
	for _, c := range selector.Conditions {
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

// findNodesDynamic finds the value for the given node using the dynamic selector
// information.
func findNodesDynamic(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	value := unwrapValue(previousValue)

	if value.Kind() == reflect.Slice {
		results := make([]*Node, 0)
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			found, err := processFindDynamicItems(selector, object)
			if err != nil {
				return nil, err
			}
			if found {
				selector.Type = "INDEX"
				selector.Index = i
				results = append(results, &Node{
					Value:    object,
					Selector: selector,
				})
			}
		}
		if len(results) > 0 {
			return results, nil
		}
		if createIfNotExists {
			selector.Type = "NEXT_AVAILABLE_INDEX"
			return []*Node{
				{
					Value:    nilValue(),
					Selector: selector,
				},
			}, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value.Kind()}
}

// findNodesAnyIndex returns a node for every value in the previous value list.
func findNodesAnyIndex(selector Selector, previousValue reflect.Value) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	value := unwrapValue(previousValue)

	if value.Kind() == reflect.Slice {
		results := make([]*Node, 0)
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			selector.Type = "INDEX"
			selector.Index = i
			results = append(results, &Node{
				Value:    object,
				Selector: selector,
			})
		}
		if len(results) > 0 {
			return results, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value.Kind()}
}

func initialiseEmptyValue(selector Selector, previousValue reflect.Value) (reflect.Value, bool) {
	switch selector.Type {
	case "PROPERTY":
		return reflect.ValueOf(map[interface{}]interface{}{}), true
	case "INDEX":
		return reflect.ValueOf([]interface{}{}), true
	case "NEXT_AVAILABLE_INDEX":
		return reflect.ValueOf([]interface{}{}), true
	case "INDEX_ANY":
		return reflect.ValueOf([]interface{}{}), true
	case "DYNAMIC":
		return reflect.ValueOf([]interface{}{}), true
	}
	return previousValue, false
}

// findNodes returns all of the nodes from the previous value that match the given selector.
func findNodes(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	initialised := false
	if createIfNotExists && !isValid(previousValue) {
		previousValue, initialised = initialiseEmptyValue(selector, previousValue)
	}

	var res []*Node
	var err error

	switch selector.Type {
	case "PROPERTY":
		res, err = findNodesProperty(selector, previousValue, createIfNotExists)
	case "INDEX":
		res, err = findNodesIndex(selector, previousValue, createIfNotExists)
	case "NEXT_AVAILABLE_INDEX":
		res, err = findNextAvailableNodes(selector, previousValue, createIfNotExists)
	case "INDEX_ANY":
		res, err = findNodesAnyIndex(selector, previousValue)
	case "DYNAMIC":
		res, err = findNodesDynamic(selector, previousValue, createIfNotExists)
	default:
		err = &UnsupportedSelector{Selector: selector.Type}
	}

	if err != nil {
		return nil, err
	}

	if initialised && res != nil && len(res) > 0 {
		for _, n := range res {
			n.wasInitialised = initialised
		}
	}

	return res, nil
}