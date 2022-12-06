package dasel

import (
	"fmt"
	"reflect"
)

// Context has scope over the entire query.
// Each individual function has its own step within the context.
// The context holds the entire data structure we're accessing/modifying.
type Context struct {
	selector          string
	selectorResolver  SelectorResolver
	steps             []*Step
	data              Value
	functions         *FunctionCollection
	createWhenMissing bool
	metadata          map[string]interface{}
}

func (c *Context) WithMetadata(key string, value interface{}) *Context {
	if c.metadata == nil {
		c.metadata = map[string]interface{}{}
	}
	c.metadata[key] = value
	return c
}

func (c *Context) Metadata(key string) interface{} {
	if c.metadata == nil {
		return nil
	}
	if val, ok := c.metadata[key]; ok {
		return val
	}
	return nil
}

func newContextWithFunctions(value interface{}, selector string, functions *FunctionCollection) *Context {
	var v Value
	if val, ok := value.(Value); ok {
		v = val
	} else {
		var reflectVal reflect.Value
		if val, ok := value.(reflect.Value); ok {
			reflectVal = val
		} else {
			reflectVal = reflect.ValueOf(value)
		}

		v = Value{
			Value: reflectVal,
		}
	}

	v.Value = makeAddressable(v.Value)

	// v.SetMapIndex(reflect.ValueOf("users"), v.MapIndex(ValueOf("users")))
	// v.MapIndex("users")

	v.setFn = func(value Value) {
		v.Unpack().Set(value.Value)
	}

	if v.Metadata("key") == nil {
		v.WithMetadata("key", "root")
	}

	return &Context{
		selector: selector,
		data:     v,
		steps: []*Step{
			{
				selector: Selector{
					funcName: "root",
					funcArgs: []string{},
				},
				index:  0,
				output: Values{v},
			},
		},
		functions:        functions,
		selectorResolver: NewSelectorResolver(selector, functions),
	}
}

func newSelectContext(value interface{}, selector string) *Context {
	return newContextWithFunctions(value, selector, standardFunctions())
}

func newPutContext(value interface{}, selector string) *Context {
	return newContextWithFunctions(value, selector, standardFunctions()).
		WithCreateWhenMissing(true)
}

func newDeleteContext(value interface{}, selector string) *Context {
	return newContextWithFunctions(value, selector, standardFunctions())
}

func derefValue(v Value) Value {
	return ValueOf(deref(v.Value))
}

func derefValues(values Values) Values {
	results := make(Values, len(values))
	for k, v := range values {
		results[k] = derefValue(v)
	}
	return results
}

// Select resolves the given selector and returns the resulting values.
func Select(root interface{}, selector string) (Values, error) {
	c := newSelectContext(root, selector)
	values, err := c.Run()
	if err != nil {
		return nil, err
	}
	return derefValues(values), nil
}

// Put resolves the given selector and writes the given value in their place.
// The root value may be changed in-place. If this is not desired you should copy the input
// value before passing it to Put.
func Put(root interface{}, selector string, value interface{}) (Value, error) {
	toSet := ValueOf(value)
	c := newPutContext(root, selector)
	values, err := c.Run()
	if err != nil {
		return Value{}, err
	}
	for _, v := range values {
		v.Set(toSet)
	}
	return c.Data(), nil
}

// Delete resolves the given selector and deletes any found values.
// The root value may be changed in-place. If this is not desired you should copy the input
// value before passing it to Delete.
func Delete(root interface{}, selector string) (Value, error) {
	c := newDeleteContext(root, selector)
	values, err := c.Run()
	if err != nil {
		return Value{}, err
	}
	for _, v := range values {
		v.Delete()
	}
	return c.Data(), nil
}

func (c *Context) subSelectContext(value interface{}, selector string) *Context {
	subC := newContextWithFunctions(value, selector, c.functions)
	subC.metadata = c.metadata
	return subC
}

func (c *Context) subSelect(value interface{}, selector string) (Values, error) {
	return c.subSelectContext(value, selector).Run()
}

// WithSelector updates c with the given selector.
func (c *Context) WithSelector(s string) *Context {
	c.selector = s
	c.selectorResolver = NewSelectorResolver(s, c.functions)
	return c
}

// WithCreateWhenMissing updates c with the given create value.
// If this value is true, elements (such as properties) will be initialised instead
// of return not found errors.
func (c *Context) WithCreateWhenMissing(create bool) *Context {
	c.createWhenMissing = create
	return c
}

// CreateWhenMissing returns true if the internal createWhenMissing value is true.
func (c *Context) CreateWhenMissing() bool {
	return c.createWhenMissing
}

// Data returns the root element of the context.
func (c *Context) Data() Value {
	return derefValue(c.data)
}

// Run calls Next repeatedly until no more steps are left.
// Returns the final Step.
func (c *Context) Run() (Values, error) {
	var res *Step
	for {
		step, err := c.Next()
		if err != nil {
			return nil, err
		}
		if step == nil {
			break
		}
		res = step
	}
	return res.output, nil
}

// Next returns the next Step, or nil if we have reached the final Selector.
func (c *Context) Next() (*Step, error) {
	nextSelector, err := c.selectorResolver.Next()
	if err != nil {
		return nil, fmt.Errorf("could not resolve selector: %w", err)
	}

	if nextSelector == nil {
		return nil, nil
	}

	nextStep := &Step{
		context:  c,
		selector: *nextSelector,
		index:    len(c.steps),
		output:   nil,
	}

	c.steps = append(c.steps, nextStep)

	if err := nextStep.execute(); err != nil {
		return nextStep, err
	}

	return nextStep, nil
}

// Step returns the step at the given index.
func (c *Context) Step(i int) *Step {
	if i < 0 || i > (len(c.steps)-1) {
		return nil
	}
	return c.steps[i]
}
