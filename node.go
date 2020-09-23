package dasel

import (
	"errors"
	"fmt"
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
	Value interface{} `json:"value"`
	// Selector is the selector for the current node.
	Selector Selector `json:"selector"`
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
		Value:    value,
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
	if n.Value == nil {
		// Set an initial value based off of the selector type.
		switch selector.Type {
		case "ROOT", "PROPERTY":
			n.Value = make(map[string]interface{})
		case "NEXT_AVAILABLE_INDEX", "INDEX", "DYNAMIC":
			n.Value = make([]interface{}, 0)
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

	previousNode.Value = newValue

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
			return nil, err
		}

		previousNode = nextNode
	}

	return previousNode, nil
}

// findValueProperty finds the value for the given node using the property selector
// information.
func findValueProperty(n *Node) (interface{}, error) {
	switch p := n.Previous.Value.(type) {
	case nil:
		return nil, &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case map[string]interface{}:
		v, ok := p[n.Selector.Property]
		if ok {
			return v, nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case map[interface{}]interface{}:
		v, ok := p[n.Selector.Property]
		if ok {
			return v, nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}
}

// propagateValueProperty does the opposite of findValueProperty.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueProperty(n *Node) error {
	switch p := n.Previous.Value.(type) {
	case nil:
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case map[string]interface{}:
		p[n.Selector.Property] = n.Value
		return nil
	case map[interface{}]interface{}:
		p[n.Selector.Property] = n.Value
		return nil
	default:
		return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}
}

// findValueIndex finds the value for the given node using the index selector
// information.
func findValueIndex(n *Node) (interface{}, error) {
	switch p := n.Previous.Value.(type) {
	case nil:
		return nil, &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		if n.Selector.Index >= 0 && n.Selector.Index < int64(len(p)) {
			return p[n.Selector.Index], nil
		}
	case []map[string]interface{}:
		if n.Selector.Index >= 0 && n.Selector.Index < int64(len(p)) {
			return p[n.Selector.Index], nil
		}
	case map[interface{}]interface{}:
		if n.Selector.Index >= 0 && n.Selector.Index < int64(len(p)) {
			return p[n.Selector.Index], nil
		}
	case []interface{}:
		if n.Selector.Index >= 0 && n.Selector.Index < int64(len(p)) {
			return p[n.Selector.Index], nil
		}
	case []string:
		if n.Selector.Index >= 0 && n.Selector.Index < int64(len(p)) {
			return p[n.Selector.Index], nil
		}
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}

	return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
}

// putNextAvailableIndex finds the value for the given node using the index selector
// information.
func putNextAvailableIndex(n *Node) (interface{}, error) {
	switch n.Previous.Value.(type) {
	case nil:
		return nil, &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		val := make(map[interface{}]interface{})
		return val, nil
	case []map[string]interface{}:
		val := make(map[string]interface{})
		return val, nil
	case []interface{}:
		val := ""
		return val, nil
	case []string:
		val := ""
		return val, nil
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}
}

// propagateValueIndex does the opposite of findValueIndex.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueIndex(n *Node) error {
	switch p := n.Previous.Value.(type) {
	case nil:
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		length := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < length {
			p[n.Selector.Index] = n.Value.(map[interface{}]interface{})
			return nil
		}
		if n.Selector.Index >= length {
			n.Previous.Value = append(p, n.Value.(map[interface{}]interface{}))
			return nil
		}
	case []map[string]interface{}:
		length := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < length {
			p[n.Selector.Index] = n.Value.(map[string]interface{})
			return nil
		}
		if n.Selector.Index >= length {
			n.Previous.Value = append(p, n.Value.(map[string]interface{}))
			return nil
		}
	case map[interface{}]interface{}:
		p[n.Selector.Index] = n.Value
		return nil
	case []interface{}:
		length := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < length {
			p[n.Selector.Index] = n.Value
			return nil
		}
		if n.Selector.Index >= length {
			n.Previous.Value = append(p, n.Value)
			return nil
		}
	case []string:
		length := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < length {
			p[n.Selector.Index] = n.Value.(string)
			return nil
		}
		if n.Selector.Index >= length {
			n.Previous.Value = append(p, n.Value.(string))
			return nil
		}
	default:
		return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}

	return &ValueNotFound{Selector: n.Selector.Current, Node: n}
}

// propagateValueNextAvailableIndex does the opposite of findValueNextAvailableIndex.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueNextAvailableIndex(n *Node) error {
	switch p := n.Previous.Value.(type) {
	case nil:
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		n.Previous.Value = append(p, n.Value.(map[interface{}]interface{}))
		return nil
	case []map[string]interface{}:
		n.Previous.Value = append(p, n.Value.(map[string]interface{}))
		return nil
	case []interface{}:
		n.Previous.Value = append(p, n.Value)
		return nil
	case []string:
		n.Previous.Value = append(p, n.Value.(string))
		return nil
	default:
		return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}
}

// processFindDynamicItem is used by findValueDynamic.
func processFindDynamicItem(n *Node, object interface{}) (interface{}, bool, error) {
	// Loop through each condition.
	allConditionsMatched := true
	for _, c := range n.Selector.Conditions {
		// If the object doesn't match any checks, return a ValueNotFound.
		found, err := c.Check(object)
		if err != nil {
			return nil, false, err
		}
		if !found {
			allConditionsMatched = false
			break
		}
	}
	if allConditionsMatched {
		return object, true, nil
	}
	return nil, false, nil
}

// findValueDynamic finds the value for the given node using the dynamic selector
// information.
func findValueDynamic(n *Node) (interface{}, error) {
	switch p := n.Previous.Value.(type) {
	case nil:
		return nil, &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		for _, object := range p {
			value, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return nil, err
			}
			if found {
				return value, nil
			}
		}
	case []map[string]interface{}:
		for _, object := range p {
			value, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return nil, err
			}
			if found {
				return value, nil
			}
		}
	case []map[string]string:
		for _, object := range p {
			value, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return nil, err
			}
			if found {
				return value, nil
			}
		}
	case []interface{}:
		for _, object := range p {
			value, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return nil, err
			}
			if found {
				return value, nil
			}
		}
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}

	return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
}

// propagateValueDynamic does the opposite of findValueDynamic.
// It finds the element in the parent the this node was created from and sets it's
// value to the value of the current node.
func propagateValueDynamic(n *Node) error {
	switch p := n.Previous.Value.(type) {
	case nil:
		return &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		for k, object := range p {
			_, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return err
			}
			if found {
				p[k] = n.Value.(map[interface{}]interface{})
				return nil
			}
		}
	case []map[string]interface{}:
		for k, object := range p {
			_, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return err
			}
			if found {
				p[k] = n.Value.(map[string]interface{})
				return nil
			}
		}
	case []map[string]string:
		for k, object := range p {
			_, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return err
			}
			if found {
				p[k] = n.Value.(map[string]string)
				return nil
			}
		}
	case []interface{}:
		for k, object := range p {
			_, found, err := processFindDynamicItem(n, object)
			if err != nil {
				return err
			}
			if found {
				p[k] = n.Value
				return nil
			}
		}
	default:
		return &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}

	return &ValueNotFound{Selector: n.Selector.Current, Node: n}
}

// findValue finds the value for the given node.
// The value is essentially pulled from the previous node, using the (already parsed) selector
// information stored on the current node.
func findValue(n *Node) (interface{}, error) {
	if n.Previous == nil {
		// previous node is required to get it's value.
		return nil, ErrMissingPreviousNode
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return findValueProperty(n)
	case "INDEX":
		return findValueIndex(n)
	case "NEXT_AVAILABLE_INDEX":
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case "DYNAMIC":
		return findValueDynamic(n)
	default:
		return nil, &UnsupportedSelector{Selector: n.Selector.Type}
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
