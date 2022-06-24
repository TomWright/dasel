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

func propagateValuePropertyWork(n *Node, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Map:
		value.SetMapIndex(reflect.ValueOf(n.Selector.Property), n.Value)
		return nil
	case reflect.Struct:
		fieldV := value.FieldByName(n.Selector.Property)
		if fieldV.IsValid() {
			fieldV.Set(n.Value)
		}
		return nil
	case reflect.Ptr:
		return propagateValuePropertyWork(n, derefValue(value))
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

// propagateValueProperty sends the value of the current node up to the previous node in the chain.
func propagateValueProperty(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}
	return propagateValuePropertyWork(n, unwrapValue(n.Previous.Value))
}

// propagateValueIndex sends the value of the current node up to the previous node in the chain.
// No need to support structs here since a struct can't have an index.
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

func deleteFromParentPropertyWork(n *Node, value reflect.Value) error {
	switch value.Kind() {
	case reflect.Map:
		value.SetMapIndex(reflect.ValueOf(n.Selector.Property), reflect.Value{})
		return nil
	case reflect.Struct:
		fieldV := value.FieldByName(n.Selector.Property)
		if fieldV.CanSet() && fieldV.IsValid() && !fieldV.IsZero() {
			fieldV.Set(reflect.New(fieldV.Type()).Elem())
		}
		return nil
	case reflect.Ptr:
		return deleteFromParentPropertyWork(n, derefValue(value))
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

// deleteFromParentProperty sends the value of the current node up to the previous node in the chain.
func deleteFromParentProperty(n *Node) error {
	if !isValid(n.Previous.Value) {
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}
	return deleteFromParentPropertyWork(n, unwrapValue(n.Previous.Value))
}

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
			value.Index(n.Selector.Index).Set(getDeletePlaceholder(value.Index(n.Selector.Index)))
		}
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// cleanupSliceDeletions scans through the given reflect and removes any invalid reflect values.
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
		if !item.IsValid() || isDeletePlaceholder(item) {
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

const deletePlaceholderKey = "dasel:delete:key"
const deletePlaceholder = "dasel:delete:me"

func getDeletePlaceholder(item reflect.Value) reflect.Value {
	switch unwrapValue(item).Kind() {
	case reflect.Map:
		return reflect.ValueOf(map[string]interface{}{
			deletePlaceholderKey: deletePlaceholder,
		})
	case reflect.Slice:
		return reflect.ValueOf([]interface{}{deletePlaceholder})
	default:
		return reflect.ValueOf(deletePlaceholder)
	}
}

func isDeletePlaceholder(item reflect.Value) bool {
	// todo : handle struct types?
	switch i := unwrapValue(item); i.Kind() {
	case reflect.Map:
		if val, ok := i.Interface().(map[string]interface{})[deletePlaceholderKey]; ok {
			if val == deletePlaceholder {
				return true
			}
		}
	case reflect.Slice:
		for _, val := range i.Interface().([]interface{}) {
			if val == deletePlaceholder {
				return true
			}
		}
	default:
		if val, ok := i.Interface().(string); ok {
			if val == deletePlaceholder {
				return true
			}
		}
	}

	return false
}
