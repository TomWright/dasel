package dasel

import (
	"fmt"
	"io"
	"reflect"
	"regexp"

	"github.com/tomwright/dasel/internal/storage"
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

// Copy returns a copy of the selector.
func (s Selector) Copy() Selector {
	return Selector{
		Raw:        s.Raw,
		Current:    s.Current,
		Remaining:  s.Remaining,
		Type:       s.Type,
		Property:   s.Property,
		Index:      s.Index,
		Conditions: s.Conditions,
	}
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
	wasInitialised bool
}

// String returns the value of the node as a string.
// No formatting is done here, you get the raw value.
func (n *Node) String() string {
	return fmt.Sprint(n.InterfaceValue())
}

// InterfaceValue returns the value stored within the node as an interface{}.
func (n *Node) InterfaceValue() interface{} {
	// We shouldn't be able to get here but this will stop a panic if we do.
	if !n.Value.IsValid() {
		return nil
	}
	return n.Value.Interface()
}

const (
	propertySelector = `(?P<property>[a-zA-Z\-_]+)`
	indexSelector    = `\[(?P<index>[0-9a-zA-Z\*]*?)\]`
)

var (
	propertyRegexp   = regexp.MustCompile(fmt.Sprintf("^\\.?%s", propertySelector))
	indexRegexp      = regexp.MustCompile(fmt.Sprintf("^\\.?%s", indexSelector))
	newDynamicRegexp = regexp.MustCompile(fmt.Sprintf("^\\.?((?:\\(.*\\))+)"))
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

// New returns a new root node with the given value.
func New(value interface{}) *Node {
	rootNode := &Node{
		Previous:     nil,
		Next:         nil,
		NextMultiple: nil,
		Selector: Selector{
			Raw:       ".",
			Current:   ".",
			Remaining: "",
			Type:      "ROOT",
			Property:  "",
		},
	}
	rootNode.setRealValue(value)
	return rootNode
}

// NewFromFile returns a new root node by parsing file using specified read parser.
func NewFromFile(path, parser string) (*Node, error) {
	readParser, err := storage.NewReadParserFromString(parser)
	if err != nil {
		return nil, err
	}

	data, err := storage.LoadFromFile(path, readParser)
	if err != nil {
		return nil, err
	}

	return New(data), nil
}

// WriteFile writes data to disk using specified write parser and options.
func (n *Node) WriteFile(writer io.Writer, parser string, compact, escapeHTML bool) error {
	writeParser, err := storage.NewWriteParserFromString(parser)
	if err != nil {
		return err
	}

	writeOptions := []storage.ReadWriteOption{
		storage.EscapeHTMLOption(escapeHTML),
	}

	if compact {
		writeOptions = append(writeOptions, storage.PrettyPrintOption(false))
	}

	value := n.InterfaceValue()
	originalValue := n.OriginalValue

	if err := storage.Write(writeParser, value, originalValue, writer, writeOptions...); err != nil {
		return fmt.Errorf("could not write to output file: %w", err)
	}

	return nil
}

func (n *Node) setValue(newValue interface{}) {
	n.Value = reflect.ValueOf(newValue)
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = newValue
	}
}

func (n *Node) setRealValue(newValue interface{}) {
	switch typed := newValue.(type) {
	case storage.RealValue:
		n.Value = reflect.ValueOf(typed.RealValue())
	default:
		n.Value = reflect.ValueOf(typed)
	}
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = newValue
	}
}

func (n *Node) setReflectValue(newValue reflect.Value) {
	n.Value = newValue
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = unwrapValue(newValue).Interface()
	}
}

func (n *Node) setRealReflectValue(newValue reflect.Value) {
	val := unwrapValue(newValue).Interface()
	switch typed := val.(type) {
	case storage.RealValue:
		n.OriginalValue = typed
		n.Value = reflect.ValueOf(typed.RealValue())
	default:
		n.Value = newValue
	}
	if n.Selector.Type == "ROOT" {
		n.OriginalValue = val
	}
}
