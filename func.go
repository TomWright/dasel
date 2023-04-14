package dasel

import (
	"fmt"
	"reflect"
	"strings"
)

type ErrUnknownFunction struct {
	Function string
}

func (e ErrUnknownFunction) Error() string {
	return fmt.Sprintf("unknown function: %s", e.Function)
}

func (e ErrUnknownFunction) Is(other error) bool {
	_, ok := other.(*ErrUnknownFunction)
	return ok
}

type ErrUnexpectedFunctionArgs struct {
	Function string
	Args     []string
	Message  string
}

func (e ErrUnexpectedFunctionArgs) Error() string {
	return fmt.Sprintf("unexpected function args: %s(%s): %s", e.Function, strings.Join(e.Args, ", "), e.Message)
}

func (e ErrUnexpectedFunctionArgs) Is(other error) bool {
	o, ok := other.(*ErrUnexpectedFunctionArgs)
	if !ok {
		return false
	}
	if o.Function != "" && o.Function != e.Function {
		return false
	}
	if o.Message != "" && o.Message != e.Message {
		return false
	}
	if o.Args != nil && !reflect.DeepEqual(o.Args, e.Args) {
		return false
	}
	return true
}

func standardFunctions() *FunctionCollection {
	collection := &FunctionCollection{}
	collection.Add(
		// Generic
		ThisFunc,
		LenFunc,
		KeyFunc,
		KeysFunc,
		MergeFunc,
		CountFunc,
		MapOfFunc,
		TypeFunc,
		JoinFunc,
		StringFunc,

		// Selectors
		IndexFunc,
		AllFunc,
		FirstFunc,
		LastFunc,
		PropertyFunc,
		AppendFunc,

		// Filters
		FilterFunc,
		FilterOrFunc,

		// Comparisons
		EqualFunc,
		MoreThanFunc,
		LessThanFunc,
		AndFunc,
		OrFunc,
		NotFunc,

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

func requireNoArgs(name string, args []string) error {
	if len(args) > 0 {
		return &ErrUnexpectedFunctionArgs{
			Function: name,
			Args:     args,
			Message:  "0 arguments expected",
		}
	}
	return nil
}

func requireExactlyXArgs(name string, args []string, x int) error {
	if len(args) != x {
		return &ErrUnexpectedFunctionArgs{
			Function: name,
			Args:     args,
			Message:  fmt.Sprintf("exactly %d arguments expected", x),
		}
	}
	return nil
}

func requireXOrMoreArgs(name string, args []string, x int) error {
	if len(args) < x {
		return &ErrUnexpectedFunctionArgs{
			Function: name,
			Args:     args,
			Message:  fmt.Sprintf("expected %d or more arguments", x),
		}
	}
	return nil
}

func requireXOrLessArgs(name string, args []string, x int) error {
	if len(args) > x {
		return &ErrUnexpectedFunctionArgs{
			Function: name,
			Args:     args,
			Message:  fmt.Sprintf("expected %d or less arguments", x),
		}
	}
	return nil
}

func requireModulusXArgs(name string, args []string, x int) error {
	if len(args)%x != 0 {
		return &ErrUnexpectedFunctionArgs{
			Function: name,
			Args:     args,
			Message:  fmt.Sprintf("expected arguments in groups of %d", x),
		}
	}
	return nil
}
