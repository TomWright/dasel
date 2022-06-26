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
	if len(n.NextMultiple) == 0 {
		return []*Node{n}
	}
	var res []*Node
	for _, nextNode := range n.NextMultiple {
		res = append(res, lastNodes(nextNode)...)
	}
	return res
}

func buildFindMultipleChain(n *Node) error {
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
	n.NextMultiple, err = findNodes(nextSelector, n, false)
	if err != nil {
		return fmt.Errorf("could not find multiple value: %w", err)
	}

	for _, next := range n.NextMultiple {
		// Add the back reference
		if next.Previous == nil {
			// This can already be set in some cases - SEARCH.
			next.Previous = n
		}

		if err := buildFindMultipleChain(next); err != nil {
			return err
		}
	}

	return nil
}

func findNodesPropertyWork(selector Selector, previousValue reflect.Value, createIfNotExists bool, value reflect.Value) ([]*Node, error) {
	switch value.Kind() {
	case reflect.Map:
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

	case reflect.Struct:
		node := &Node{
			Value:    nilValue(),
			Selector: selector,
		}
		fieldV := value.FieldByName(selector.Property)
		if fieldV.IsValid() {
			node.Value = fieldV
			return []*Node{node}, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}

	case reflect.Ptr:
		return findNodesPropertyWork(selector, previousValue, createIfNotExists, derefValue(value))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNodesProperty finds the value for the given node using the property selector
// information.
func findNodesProperty(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}

	if selector.Property == "-" {
		res, err := findNodesPropertyKeys(selector, previousValue, createIfNotExists)
		return res, err
	}

	return findNodesPropertyWork(selector, previousValue, createIfNotExists, unwrapValue(previousValue))
}

func findNodesPropertyKeysWork(selector Selector, previousValue reflect.Value, createIfNotExists bool, value reflect.Value) ([]*Node, error) {
	results := make([]*Node, 0)

	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			sel := selector.Copy()
			sel.Type = "INDEX"
			sel.Index = i
			results = append(results, &Node{
				Value:    reflect.ValueOf(fmt.Sprint(i)),
				Selector: sel,
			})
		}
		return results, nil

	case reflect.Map:
		for _, key := range value.MapKeys() {
			sel := selector.Copy()
			sel.Type = "PROPERTY"
			sel.Property = key.String()
			results = append(results, &Node{
				Value:    reflect.ValueOf(key.String()),
				Selector: sel,
			})
		}
		return results, nil

	case reflect.Struct:
		valueType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := valueType.Field(i)
			sel := selector.Copy()
			sel.Type = "PROPERTY"
			sel.Property = field.Name
			results = append(results, &Node{
				Value:    reflect.ValueOf(field.Name),
				Selector: sel,
			})
		}
		return results, nil

	case reflect.Ptr:
		return findNodesPropertyKeysWork(selector, previousValue, createIfNotExists, derefValue(value))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

func findNodesPropertyKeys(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	if createIfNotExists {
		return nil, &UnsupportedSelector{Selector: selector.Raw}
	}
	return findNodesPropertyKeysWork(selector, previousValue, createIfNotExists, unwrapValue(previousValue))
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

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNextAvailableIndexNodes finds the value for the given node using the index selector
// information.
func findNextAvailableIndexNodes(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
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
func processFindDynamicItems(selector Selector, object reflect.Value, key string) (bool, error) {
	// Loop through each condition.
	allConditionsMatched := true
	for _, c := range selector.Conditions {
		// If the object doesn't match any checks, return a ValueNotFound.

		var found bool
		var err error
		switch cond := c.(type) {
		case *KeyEqualCondition:
			found, err = cond.Check(reflect.ValueOf(key))
		default:
			found, err = cond.Check(object)
		}

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

func findNodesDynamicWork(selector Selector, previousValue reflect.Value, createIfNotExists bool, value reflect.Value) ([]*Node, error) {
	switch value.Kind() {
	case reflect.Slice:
		results := make([]*Node, 0)
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			found, err := processFindDynamicItems(selector, object, fmt.Sprint(i))
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

	case reflect.Map:
		results := make([]*Node, 0)
		for _, key := range value.MapKeys() {
			object := value.MapIndex(key)
			found, err := processFindDynamicItems(selector, object, key.String())
			if err != nil {
				return nil, err
			}
			if found {
				selector.Type = "PROPERTY"
				selector.Property = key.String()
				results = append(results, &Node{
					Value:    object,
					Selector: selector,
				})
			}
		}
		if len(results) > 0 {
			return results, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}

	case reflect.Struct:
		results := make([]*Node, 0)
		valueType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := valueType.Field(i)
			object := value.Field(i)
			found, err := processFindDynamicItems(selector, object, field.Name)
			if err != nil {
				return nil, err
			}
			if found {
				selector.Type = "PROPERTY"
				selector.Property = field.Name
				results = append(results, &Node{
					Value:    object,
					Selector: selector,
				})
			}
		}
		if len(results) > 0 {
			return results, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}

	case reflect.Ptr:
		return findNodesDynamicWork(selector, previousValue, createIfNotExists, derefValue(value))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNodesDynamic finds the value for the given node using the dynamic selector
// information.
func findNodesDynamic(selector Selector, previousValue reflect.Value, createIfNotExists bool) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	return findNodesDynamicWork(selector, previousValue, createIfNotExists, unwrapValue(previousValue))
}

func findNodesSearchRecursiveSubNode(selector Selector, subNode *Node, key string, createIfNotExists bool) ([]*Node, error) {
	subResults, err := findNodesSearchRecursive(selector, subNode, createIfNotExists, false)
	if err != nil {
		return nil, fmt.Errorf("could not find nodes search recursive: %w", err)
	}

	// Loop through each condition.
	allConditionsMatched := true
sliceConditionLoop:
	for _, c := range selector.Conditions {
		var found bool
		var err error

		switch cond := c.(type) {
		case *KeyEqualCondition:
			found, err = cond.Check(reflect.ValueOf(key))
		default:
			found, err = cond.Check(subNode.Value)
		}
		if err != nil || !found {
			allConditionsMatched = false
			break sliceConditionLoop
		}
	}

	results := make([]*Node, 0)

	if allConditionsMatched {
		results = append(results, subNode)
	}
	if len(subResults) > 0 {
		results = append(results, subResults...)
	}

	return results, nil
}

func findNodesSearchRecursiveWork(selector Selector, previousNode *Node, createIfNotExists bool, firstNode bool, value reflect.Value) ([]*Node, error) {
	results := make([]*Node, 0)

	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)

			subNode := &Node{
				Previous: previousNode,
				Value:    object,
				Selector: selector.Copy(),
			}
			subNode.Selector.Type = "INDEX"
			subNode.Selector.Index = i

			newResults, err := findNodesSearchRecursiveSubNode(selector, subNode, fmt.Sprint(subNode.Selector.Index), createIfNotExists)
			if err != nil {
				return nil, err
			}
			results = append(results, newResults...)
		}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			object := value.MapIndex(key)

			subNode := &Node{
				Previous: previousNode,
				Value:    object,
				Selector: selector.Copy(),
			}
			subNode.Selector.Type = "PROPERTY"
			subNode.Selector.Property = fmt.Sprint(key.Interface())

			newResults, err := findNodesSearchRecursiveSubNode(selector, subNode, subNode.Selector.Property, createIfNotExists)
			if err != nil {
				return nil, err
			}
			results = append(results, newResults...)
		}

	case reflect.Struct:
		valueType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := valueType.Field(i)
			object := value.Field(i)

			subNode := &Node{
				Previous: previousNode,
				Value:    object,
				Selector: selector.Copy(),
			}
			subNode.Selector.Type = "PROPERTY"
			subNode.Selector.Property = field.Name

			newResults, err := findNodesSearchRecursiveSubNode(selector, subNode, subNode.Selector.Property, createIfNotExists)
			if err != nil {
				return nil, err
			}
			results = append(results, newResults...)
		}

	case reflect.Ptr:
		return findNodesSearchRecursiveWork(selector, previousNode, createIfNotExists, firstNode, derefValue(value))
	}

	return results, nil
}

// findNodesSearchRecursive iterates through the value of the previous node and creates a new node for each element.
// If any of those nodes match the checks they are returned.
func findNodesSearchRecursive(selector Selector, previousNode *Node, createIfNotExists bool, firstNode bool) ([]*Node, error) {
	if !isValid(previousNode.Value) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	return findNodesSearchRecursiveWork(selector, previousNode, createIfNotExists, firstNode, unwrapValue(previousNode.Value))
}

// findNodesSearch finds all available nodes by recursively searching the previous value.
func findNodesSearch(selector Selector, previousNode *Node, createIfNotExists bool) ([]*Node, error) {
	res, err := findNodesSearchRecursive(selector, previousNode, createIfNotExists, true)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, &ValueNotFound{
			Selector:      selector.Current,
			PreviousValue: previousNode.Value,
		}
	}
	return res, nil
}

func findNodesAnyIndexWork(selector Selector, previousValue reflect.Value, value reflect.Value) ([]*Node, error) {
	switch value.Kind() {
	case reflect.Slice:
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

	case reflect.Map:
		results := make([]*Node, 0)
		for _, key := range value.MapKeys() {
			object := value.MapIndex(key)
			selector.Type = "PROPERTY"
			selector.Property = key.String()
			results = append(results, &Node{
				Value:    object,
				Selector: selector,
			})
		}
		if len(results) > 0 {
			return results, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}

	case reflect.Struct:
		valueType := value.Type()
		results := make([]*Node, 0)
		for i := 0; i < valueType.NumField(); i++ {
			field := valueType.Field(i)
			object := value.Field(i)
			selector.Type = "PROPERTY"
			selector.Property = field.Name
			results = append(results, &Node{
				Value:    object,
				Selector: selector,
			})
		}
		if len(results) > 0 {
			return results, nil
		}
		return nil, &ValueNotFound{Selector: selector.Current, PreviousValue: previousValue}

	case reflect.Ptr:
		return findNodesAnyIndexWork(selector, previousValue, derefValue(previousValue))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNodesAnyIndex returns a node for every value in the previous value list.
func findNodesAnyIndex(selector Selector, previousValue reflect.Value) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	return findNodesAnyIndexWork(selector, previousValue, unwrapValue(previousValue))
}

func findNodesLengthWork(selector Selector, previousValue reflect.Value, value reflect.Value) ([]*Node, error) {
	switch value.Kind() {
	case reflect.Slice:
		node := &Node{
			Value:    reflect.ValueOf(value.Len()),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Map:
		node := &Node{
			Value:    reflect.ValueOf(value.Len()),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.String:
		node := &Node{
			Value:    reflect.ValueOf(value.Len()),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Struct:
		node := &Node{
			Value:    reflect.ValueOf(value.NumField()),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Ptr:
		return findNodesLengthWork(selector, previousValue, derefValue(value))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNodesLength returns the length
func findNodesLength(selector Selector, previousValue reflect.Value) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	return findNodesLengthWork(selector, previousValue, unwrapValue(previousValue))
}

func findNodesTypeWork(selector Selector, value reflect.Value) ([]*Node, error) {
	switch value.Kind() {
	case reflect.Slice:
		node := &Node{
			Value:    reflect.ValueOf("array"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Map:
		node := &Node{
			Value:    reflect.ValueOf("map"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.String:
		node := &Node{
			Value:    reflect.ValueOf("string"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		node := &Node{
			Value:    reflect.ValueOf("int"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Float32, reflect.Float64:
		node := &Node{
			Value:    reflect.ValueOf("float"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Bool:
		node := &Node{
			Value:    reflect.ValueOf("bool"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Struct:
		node := &Node{
			Value:    reflect.ValueOf("struct"),
			Selector: selector,
		}
		return []*Node{node}, nil

	case reflect.Ptr:
		return findNodesTypeWork(selector, derefValue(value))
	}

	return nil, &UnsupportedTypeForSelector{Selector: selector, Value: value}
}

// findNodesType returns the type
func findNodesType(selector Selector, previousValue reflect.Value) ([]*Node, error) {
	if !isValid(previousValue) {
		return nil, &UnexpectedPreviousNilValue{Selector: selector.Raw}
	}
	return findNodesTypeWork(selector, unwrapValue(previousValue))
}

func initialiseEmptyValue(selector Selector, previousValue reflect.Value) reflect.Value {
	switch selector.Type {
	case "PROPERTY":
		return reflect.ValueOf(map[string]interface{}{})
	case "INDEX", "NEXT_AVAILABLE_INDEX", "INDEX_ANY", "DYNAMIC":
		return reflect.ValueOf([]interface{}{})
	}
	return previousValue
}

// findNodes returns all of the nodes from the previous value that match the given selector.
func findNodes(selector Selector, previousNode *Node, createIfNotExists bool) ([]*Node, error) {
	if createIfNotExists && !isValid(previousNode.Value) {
		previousNode.Value = initialiseEmptyValue(selector, previousNode.Value)
	}

	var res []*Node
	var err error

	switch selector.Type {
	case "PROPERTY":
		res, err = findNodesProperty(selector, previousNode.Value, createIfNotExists)
	case "INDEX":
		res, err = findNodesIndex(selector, previousNode.Value, createIfNotExists)
	case "NEXT_AVAILABLE_INDEX":
		res, err = findNextAvailableIndexNodes(selector, previousNode.Value, createIfNotExists)
	case "INDEX_ANY":
		res, err = findNodesAnyIndex(selector, previousNode.Value)
	case "LENGTH":
		res, err = findNodesLength(selector, previousNode.Value)
	case "TYPE":
		res, err = findNodesType(selector, previousNode.Value)
	case "DYNAMIC":
		res, err = findNodesDynamic(selector, previousNode.Value, createIfNotExists)
	case "SEARCH":
		res, err = findNodesSearch(selector, previousNode, createIfNotExists)
	default:
		err = &UnsupportedSelector{Selector: selector.Raw}
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}
