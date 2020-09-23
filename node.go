package dasel

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Selector represents the selector for a node.
type Selector struct {
	// Raw is the full selector.
	Raw string `json:"raw"`
	// Current is the selector to be used with the current node.
	Current string `json:"current"`
	// Remaining is the remaining parts of the Raw selector.
	Remaining string `json:"remaining"`
	// Type is the type of the selector.
	Type string `json:"type"`
	// Property is the name of the property this selector targets, if applicable.
	Property string `json:"property,omitempty"`
	// Index is the index to use if applicable.
	Index int64 `json:"index,omitempty"`
	// Conditions contains a set of conditions to optionally match a target.
	Conditions []Condition `json:"conditions,omitempty"`
}

// Node represents a single node in the chain of nodes for a selector.
type Node struct {
	// Previous is the previous node in the chain.
	Previous *Node `json:"-"`
	// Next is the next node in the chain.
	Next *Node `json:"next,omitempty"`

	// Value is the value of the current node.
	Value reflect.Value `json:"value"`
	// Selector is the selector for the current node.
	Selector Selector `json:"selector"`
}

func (n *Node) InterfaceValue() interface{} {
	return n.Value.Interface()
}

const (
	propertySelector = `(?P<property>[a-zA-Z\-_]+)`
	indexSelector    = `\[(?P<index>[0-9a-zA-Z]*?)\]`
	dynamicSelector  = `\((?P<name>[a-zA-Z\-_]+)(?P<comparison>=|<|>)(?P<name>.*?)\)`
)

var (
	// firstNodeRegexp = regexp.MustCompile(`(.+\.?)*?`)
	propertyRegexp        = regexp.MustCompile(fmt.Sprintf("^\\.?%s", propertySelector))
	indexRegexp           = regexp.MustCompile(fmt.Sprintf("^\\.?%s", indexSelector))
	dynamicRegexp         = regexp.MustCompile(fmt.Sprintf("^\\.?(?:%s)+", dynamicSelector))
	multipleDynamicRegexp = regexp.MustCompile(fmt.Sprintf("%s+", dynamicSelector))
)

func isValid(value reflect.Value) bool {
	return value.IsValid() && !safeIsNil(value)
}

func safeIsNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return value.IsNil()
	}
	return false
}

func nilValue() reflect.Value {
	return reflect.ValueOf(nil)
}

func unwrapValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Interface {
		return value.Elem()
	}
	return value
}

// ParseSelector parses the given selector string and returns a Selector.
func ParseSelector(selector string) (Selector, error) {
	sel := Selector{
		Raw:        selector,
		Current:    "",
		Remaining:  "",
		Type:       "",
		Property:   "",
		Conditions: make([]Condition, 0),
	}

	if match := propertyRegexp.FindStringSubmatch(selector); len(match) != 0 {
		sel.Type = "PROPERTY"
		sel.Current = match[0]
		sel.Property = match[1]
	} else if match := indexRegexp.FindStringSubmatch(selector); len(match) != 0 {
		sel.Current = match[0]
		if match[1] == "" {
			sel.Type = "NEXT_AVAILABLE_INDEX"
		} else {
			sel.Type = "INDEX"
			var err error
			sel.Index, err = strconv.ParseInt(match[1], 10, 64)
			if err != nil {
				return sel, &InvalidIndexErr{Index: match[1]}
			}
		}
	} else if match := dynamicRegexp.FindString(selector); match != "" {
		sel.Current = match
		matches := multipleDynamicRegexp.FindAllStringSubmatch(match, -1)
		for _, m := range matches {
			var cond Condition
			switch m[2] {
			case "=":
				cond = &EqualCondition{
					Key:   m[1],
					Value: m[3],
				}
			default:
				return sel, &UnknownComparisonOperatorErr{Operator: m[2]}
			}

			sel.Conditions = append(sel.Conditions, cond)
		}
		sel.Type = "DYNAMIC"
	}

	sel.Remaining = strings.TrimPrefix(sel.Raw, sel.Current)

	return sel, nil
}

// New returns a new root note with the given value.
func New(value interface{}) *Node {
	rootNode := &Node{
		Previous: nil,
		Next:     nil,
		Value:    reflect.ValueOf(value),
		Selector: Selector{
			Raw:       ".",
			Current:   ".",
			Remaining: "",
			Type:      "ROOT",
			Property:  "",
		},
	}
	return rootNode
}

// Initialise sets up the node to ensure it has a value.
func (n *Node) Initialise(selector Selector) {
	if !isValid(n.Value) {
		// Set an initial value based off of the selector type.
		switch selector.Type {
		case "ROOT", "PROPERTY":
			n.Value = reflect.ValueOf(make(map[string]interface{}))
		case "NEXT_AVAILABLE_INDEX", "INDEX", "DYNAMIC":
			n.Value = reflect.ValueOf(make([]interface{}, 0))
		}
	}
}

// Put finds the node using the given selector and updates it's value.
// It then attempts to propagate the value back up the chain to the root element.
func (n *Node) Put(selector string, newValue interface{}) error {
	n.Selector.Remaining = selector
	rootNode := n
	previousNode := rootNode
	var nextNode *Node
	var err error

	for {
		if previousNode.Selector.Remaining == "" {
			break
		}

		nextNode = &Node{}

		// Parse the selector.
		nextNode.Selector, err = ParseSelector(previousNode.Selector.Remaining)
		if err != nil {
			return fmt.Errorf("failed to parse selector: %w", err)
		}

		previousNode.Initialise(nextNode.Selector)

		// Link the nodes.
		previousNode.Next = nextNode
		nextNode.Previous = previousNode

		// Populate the value for the new node.
		nextNode.Value, err = findValue(nextNode)
		if err != nil {
			var valueNotFoundErr *ValueNotFound
			if errors.As(err, &valueNotFoundErr) {
				if nextNode.Selector.Type == "NEXT_AVAILABLE_INDEX" {
					nextNode.Value, err = putNextAvailableIndex(nextNode)
					if err != nil {
						return fmt.Errorf("could not put next available index: %w", err)
					}
				}
			} else {
				return fmt.Errorf("could not find value: %w", err)
			}
		}

		previousNode = nextNode
	}

	previousNode.Value = reflect.ValueOf(newValue)

	for {
		if previousNode.Previous == nil {
			break
		}
		if err := propagateValue(previousNode); err != nil {
			return fmt.Errorf("could not propagate value: %w", err)
		}
		previousNode = previousNode.Previous
	}

	return nil
}

// Query uses the given selector to query the current node and return the result.
func (n *Node) Query(selector string) (*Node, error) {
	n.Selector.Remaining = selector
	rootNode := n
	previousNode := rootNode
	var nextNode *Node
	var err error

	for {
		if previousNode == nil || previousNode.Selector.Remaining == "" {
			break
		}

		nextNode = &Node{}

		// Parse the selector.
		nextNode.Selector, err = ParseSelector(previousNode.Selector.Remaining)
		if err != nil {
			return nil, fmt.Errorf("failed to parse selector: %w", err)
		}

		// Link the nodes.
		previousNode.Next = nextNode
		nextNode.Previous = previousNode

		nextNode.Value, err = findValue(nextNode)
		// Populate the value for the new node.
		if err != nil {
			return nil, fmt.Errorf("could not find value: %w", err)
		}

		previousNode = nextNode
	}

	return previousNode, nil
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
		return nil
	}

	return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
}

// findValueIndex finds the value for the given node using the index selector
// information.
func findValueIndex(n *Node) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		if n.Selector.Index >= 0 && n.Selector.Index < int64(value.Len()) {
			return value.Index(int(n.Selector.Index)), nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
}

// putNextAvailableIndex finds the value for the given node using the index selector
// information.
func putNextAvailableIndex(n *Node) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		if value.Len() == 0 {
			return reflect.ValueOf(map[interface{}]interface{}{}), nil
		}
		return reflect.New(value.Index(0).Type()).Elem(), nil
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value.Kind()}
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
		if n.Selector.Index >= 0 && n.Selector.Index < int64(value.Len()) {
			value.Index(int(n.Selector.Index)).Set(n.Value)
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

// propagateValueDynamic does the opposite of findValueDynamic.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
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
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case "DYNAMIC":
		return findValueDynamic(n)
	default:
		return nilValue(), &UnsupportedSelector{Selector: n.Selector.Type}
	}
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
