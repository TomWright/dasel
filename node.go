package dasel

import (
	"fmt"
	"github.com/tomwright/dasel/internal/storage"
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
	Index int `json:"index,omitempty"`
	// Conditions contains a set of conditions to optionally match a target.
	Conditions []Condition `json:"conditions,omitempty"`
}

// Node represents a single node in the chain of nodes for a selector.
type Node struct {
	// Previous is the previous node in the chain.
	Previous *Node `json:"-"`
	// Next contains the next node in the chain.
	// This is used with Query and Put requests.
	Next *Node `json:"next,omitempty"`
	// NextMultiple contains the next nodes in the chain.
	// This is used with QueryMultiple and PutMultiple requests.
	// When a major version change occurs this will completely replace Next.
	NextMultiple []*Node `json:"nextMultiple,omitempty"`
	// OriginalValue is the value returned from the parser.
	// In most cases this is the same as Value, but is different for thr YAML parser
	// as it contains information on the original document.
	OriginalValue interface{} `json:"-"`
	// Value is the value of the current node.
	Value reflect.Value `json:"value"`
	// Selector is the selector for the current node.
	Selector       Selector `json:"selector"`
	propagated     bool
	wasInitialised bool
}

// InterfaceValue returns the value stored within the node as an interface{}.
func (n *Node) InterfaceValue() interface{} {
	return n.Value.Interface()
}

const (
	propertySelector = `(?P<property>[a-zA-Z\-_]+)`
	indexSelector    = `\[(?P<index>[0-9a-zA-Z\*]*?)\]`
	dynamicSelector  = `(?P<name>.+)(?P<comparison>=|<|>)(?P<value>.+)`
)

var (
	propertyRegexp        = regexp.MustCompile(fmt.Sprintf("^\\.?%s", propertySelector))
	indexRegexp           = regexp.MustCompile(fmt.Sprintf("^\\.?%s", indexSelector))
	dynamicSelectorRegexp = regexp.MustCompile(fmt.Sprintf("%s", dynamicSelector))
	newDynamicRegexp      = regexp.MustCompile(fmt.Sprintf("^\\.?((?:\\(.*\\))+)"))
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
	// value = reflect.Indirect(value)
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
		switch match[1] {
		case "":
			sel.Type = "NEXT_AVAILABLE_INDEX"
		case "*":
			sel.Type = "INDEX_ANY"
		default:
			sel.Type = "INDEX"
			var err error
			index, err := strconv.ParseInt(match[1], 10, 32)
			if err != nil {
				return sel, &InvalidIndexErr{Index: match[1]}
			}
			sel.Index = int(index)
		}
	} else {
		// todo : re-work this logic to base the entire parsing using ExtractNextSelector.
		// This will be much easier instead of regex.
		nextSelectorString := ExtractNextSelector(selector)

		// Check if the selector starts with an open bracket.
		if !strings.HasPrefix(strings.TrimPrefix(nextSelectorString, "."), "(") {
			return sel, &UnsupportedSelector{Selector: nextSelectorString}
		}

		sel.Current = nextSelectorString

		dynamicGroups, err := DynamicSelectorToGroups(nextSelectorString)
		if err != nil {
			return sel, err
		}

		for _, g := range dynamicGroups {
			m := dynamicSelectorRegexp.FindStringSubmatch(g)

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

// New returns a new root node with the given value.
func New(value interface{}) *Node {
	var baseValue reflect.Value
	switch typed := value.(type) {
	case storage.RealValue:
		baseValue = reflect.ValueOf(typed.RealValue())
	default:
		baseValue = reflect.ValueOf(value)
	}
	rootNode := &Node{
		Previous:      nil,
		Next:          nil,
		NextMultiple:  nil,
		OriginalValue: value,
		Value:         baseValue,
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
