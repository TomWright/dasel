package query

import (
	"fmt"
	"reflect"
)

// Step is a single step in the query.
// Each function call has its own step.
// Each value in the output is simply a pointer to the actual data point in the context data.
type Step struct {
	selector Selector
	index    int
	output   Values
}

func (s *Step) Selector() Selector {
	return s.selector
}

func (s *Step) Index() int {
	return s.index
}

func (s *Step) Output() Values {
	return s.output
}

// Context has scope over the entire query.
// Each individual function has its own step within the context.
// The context holds the entire data structure we're accessing/modifying.
type Context struct {
	selector         string
	selectorResolver SelectorResolver
	steps            []*Step
	data             Value
	functions        *FunctionCollection
}

func newContextWithFunctions(value interface{}, selector string, functions *FunctionCollection) *Context {
	var v Value
	if val, ok := value.(Value); ok {
		v = val
	} else {
		reflectVal := reflect.ValueOf(value)
		v = Value{
			Value: reflectVal,
			setFn: func(value Value) {
				reflectVal.Set(value.Value)
			},
			metadata: map[string]interface{}{},
		}
		v.metadata = map[string]interface{}{
			"type": v.Unpack().Type().String(),
			"key":  "root",
		}
	}

	if v.Kind() == reflect.Struct && v.Type().Kind() != reflect.Ptr {
		panic("new context must be passed a pointer")
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

// NewContext returns a new query context.
func NewContext(value interface{}, selector string) *Context {
	return newContextWithFunctions(value, selector, standardFunctions())
}

// Run calls Next repeatedly until no more steps are left.
// Returns the final Step.
func (c *Context) Run() (*Step, error) {
	var res *Step
	for {
		step, err := c.Next()
		if err != nil {
			return res, err
		}
		if step == nil {
			break
		}
		res = step
	}
	return res, nil
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
		selector: *nextSelector,
		index:    len(c.steps),
		output:   nil,
	}

	c.steps = append(c.steps, nextStep)

	if err := c.processStep(nextStep); err != nil {
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

func (c *Context) processStep(step *Step) error {
	f, err := c.functions.Get(step.selector.funcName)
	if err != nil {
		return err
	}
	output, err := f(c, step, step.selector.funcArgs)
	step.output = output
	return err
}

func (c *Context) inputValue(s *Step) Values {
	prevStep := c.Step(s.index - 1)
	if prevStep == nil {
		return Values{}
	}
	return prevStep.output
}

func (c *Context) subContext(value interface{}, selector string) *Context {
	return newContextWithFunctions(value, selector, c.functions)
}

func performSubQuery(c *Context, value Value, selector string) (Values, error) {
	matchC := c.subContext(value, selector)
	finalStep, err := matchC.Run()
	if err != nil {
		return nil, err
	}
	return finalStep.output, nil
}
