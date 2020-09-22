package dasel

import (
	"encoding/json"
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

// String returns a string representation of the node. It does this by marshaling it.
func (n *Node) String() string {
	b, err := json.MarshalIndent(n, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

const (
	propertySelector = `(?P<property>[a-zA-Z\-_]+)`
	indexSelector    = `\[(?P<index>[0-9]*?)\]`
	dynamicSelector  = `\((?P<name>[a-zA-Z\-_]+)(?P<comparison>=)(?P<name>.*?)\)`
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

// Query uses the given selector to query the current node and return the result.
func (n Node) Query(selector string) (*Node, error) {
	n.Selector.Remaining = selector
	rootNode := &n
	var previousNode = rootNode
	var nextNode *Node
	var err error

	for {
		if previousNode == nil {
			break
		}
		if previousNode.Selector.Remaining == "" {
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

		nextNode.Value, err = FindValue(nextNode)
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

// findValueIndex finds the value for the given node using the index selector
// information.
func findValueIndex(n *Node) (interface{}, error) {
	switch p := n.Previous.Value.(type) {
	case nil:
		return nil, &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	case []map[interface{}]interface{}:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case []map[string]interface{}:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case map[interface{}]interface{}:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case map[int]interface{}:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[int(n.Selector.Index)], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case []interface{}:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case []string:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	case []int:
		l := int64(len(p))
		if n.Selector.Index >= 0 && n.Selector.Index < l {
			return p[n.Selector.Index], nil
		}
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
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
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
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
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
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
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
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
		return nil, &ValueNotFound{Selector: n.Selector.Current, Node: n}
	default:
		return nil, &UnsupportedTypeForSelector{Selector: n.Selector, Value: n.Previous.Value}
	}
}

// FindValue finds the value for the given node.
// The value is essentially pulled from the previous node, using the (already parsed) selector
// information stored on the current node.
func FindValue(n *Node) (interface{}, error) {
	if n.Previous == nil {
		// previous node is required to get it's value.
		return nil, ErrMissingPreviousNode
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return findValueProperty(n)
	case "INDEX":
		return findValueIndex(n)
	case "DYNAMIC":
		return findValueDynamic(n)
	default:
		return nil, &UnsupportedSelector{Selector: n.Selector.Type}
	}
}
