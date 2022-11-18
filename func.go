package dasel

import "fmt"

type ErrUnknownFunction struct {
	Function string
}

func (e ErrUnknownFunction) Error() string {
	return fmt.Sprintf("unknown function: %s", e.Function)
}

func (e ErrUnknownFunction) Is(other error) bool {
	_, ok := other.(ErrUnknownFunction)
	return ok
}

func standardFunctions() *FunctionCollection {
	collection := &FunctionCollection{}
	collection.Add(
		// Generic
		ThisFunc,
		LenFunc,
		KeyFunc,

		// Selectors
		IndexFunc,
		AllFunc,
		FirstFunc,
		LastFunc,
		PropertyFunc,

		// Filters
		FilterFunc,
		FilterOrFunc,

		// Comparisons
		EqualFunc,
		MoreThanFunc,
		LessThanFunc,

		// Metadata
		MetadataFunc,
		ParentFunc,
	)
	return collection
}

// SelectorFunc is a function that can be executed in a selector.
type SelectorFunc func(c *Context, step *Step, args []string) (Values, error)

type FunctionCollection struct {
	functions []Function
}

func (fc *FunctionCollection) ParseSelector(part string) *Selector {
	for _, f := range fc.functions {
		if s := f.AlternativeSelector(part); s != nil {
			return s
		}
	}
	return nil
}

func (fc *FunctionCollection) Add(fs ...Function) {
	fc.functions = append(fc.functions, fs...)
}

func (fc *FunctionCollection) GetAll() map[string]SelectorFunc {
	res := make(map[string]SelectorFunc)
	for _, f := range fc.functions {
		res[f.Name()] = f.Run
	}
	return res
}

func (fc *FunctionCollection) Get(name string) (SelectorFunc, error) {
	if f, ok := fc.GetAll()[name]; ok {
		return f, nil
	}
	return nil, &ErrUnknownFunction{Function: name}
}

type Function interface {
	Name() string
	Run(c *Context, s *Step, args []string) (Values, error)
	AlternativeSelector(part string) *Selector
}

type BasicFunction struct {
	name                  string
	runFn                 func(c *Context, s *Step, args []string) (Values, error)
	alternativeSelectorFn func(part string) *Selector
}

func (bf BasicFunction) Name() string {
	return bf.name
}

func (bf BasicFunction) Run(c *Context, s *Step, args []string) (Values, error) {
	return bf.runFn(c, s, args)
}

func (bf BasicFunction) AlternativeSelector(part string) *Selector {
	if bf.alternativeSelectorFn == nil {
		return nil
	}
	return bf.alternativeSelectorFn(part)
}
