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
			Value:    reflectVal,
			metadata: map[string]interface{}{},
		}
	}

	// Make sure we have an addressable root value.
	if !v.CanAddr() {
		pointerValue := reflect.New(v.Value.Type())
		pointerValue.Elem().Set(v.Value)
		v.Value = pointerValue
	}

	v.setFn = func(value Value) {
		v.Unpack().Set(value.Value)
	}

	if v.metadata == nil {
		v.metadata = map[string]interface{}{}
	}

	if v.Metadata("key") == nil {
		v.WithMetadata("key", "root")
	}
	v.WithMetadata("type", v.Unpack().Type().String())

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

func Select(root interface{}, selector string) (Values, error) {
	c := newSelectContext(root, selector)
	return c.Run()
}

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
	return newContextWithFunctions(value, selector, c.functions)
}

func (c *Context) subSelect(value interface{}, selector string) (Values, error) {
	return c.subSelectContext(value, selector).Run()
}

func (c *Context) WithSelector(s string) *Context {
	c.selector = s
	c.selectorResolver = NewSelectorResolver(s, c.functions)
	return c
}

func (c *Context) WithCreateWhenMissing(create bool) *Context {
	c.createWhenMissing = create
	return c
}

func (c *Context) CreateWhenMissing() bool {
	return c.createWhenMissing
}

func (c *Context) Data() Value {
	return c.data
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
