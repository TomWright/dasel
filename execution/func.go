package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

var (
	// DefaultFuncCollection is the default collection of functions that can be executed.
	DefaultFuncCollection = NewFuncCollection(
		FuncLen,
		FuncAdd,
		FuncToString,
		FuncToInt,
		FuncToFloat,
		FuncMerge,
		FuncReverse,
		FuncTypeOf,
		FuncMax,
		FuncMin,
		FuncIgnore,
		FuncBase64Encode,
		FuncBase64Decode,
		FuncParse,
		FuncReadFile,
		FuncHas,
		FuncGet,
		FuncContains,
		FuncSum,
		FuncJoin,
		FuncReplace,
	)
)

// ArgsValidator is a function that validates the arguments passed to a function.
type ArgsValidator func(ctx context.Context, name string, args model.Values) error

// ValidateArgsExactly returns an ArgsValidator that validates that the number of arguments passed to a function is exactly the expected number.
func ValidateArgsExactly(expected int) ArgsValidator {
	return func(ctx context.Context, name string, args model.Values) error {
		if len(args) == expected {
			return nil
		}
		return fmt.Errorf("func %q expects exactly %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMin returns an ArgsValidator that validates that the number of arguments passed to a function is at least the expected number.
func ValidateArgsMin(expected int) ArgsValidator {
	return func(ctx context.Context, name string, args model.Values) error {
		if len(args) >= expected {
			return nil
		}
		return fmt.Errorf("func %q expects at least %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMax returns an ArgsValidator that validates that the number of arguments passed to a function is at most the expected number.
func ValidateArgsMax(expected int) ArgsValidator {
	return func(ctx context.Context, name string, args model.Values) error {
		if len(args) <= expected {
			return nil
		}
		return fmt.Errorf("func %q expects no more than %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMinMax returns an ArgsValidator that validates that the number of arguments passed to a function is between the min and max expected numbers.
func ValidateArgsMinMax(min int, max int) ArgsValidator {
	return func(ctx context.Context, name string, args model.Values) error {
		if len(args) >= min && len(args) <= max {
			return nil
		}
		return fmt.Errorf("func %q expects between %d and %d arguments, got %d", name, min, max, len(args))
	}
}

// Func represents a function that can be executed.
type Func struct {
	name          string
	handler       FuncFn
	argsValidator ArgsValidator
}

// Handler returns a FuncFn that can be used to execute the function.
func (f *Func) Handler() FuncFn {
	return func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		if f.argsValidator != nil {
			if err := f.argsValidator(ctx, f.name, args); err != nil {
				return nil, err
			}
		}
		res, err := f.handler(ctx, data, args)
		if err != nil {
			return nil, fmt.Errorf("error execution func %q: %w", f.name, err)
		}
		return res, nil
	}
}

// NewFunc creates a new Func.
func NewFunc(name string, handler FuncFn, argsValidator ArgsValidator) *Func {
	return &Func{
		name:          name,
		handler:       handler,
		argsValidator: argsValidator,
	}
}

// FuncFn is a function that can be executed.
type FuncFn func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error)

// FuncCollection is a collection of functions that can be executed.
type FuncCollection map[string]FuncFn

// NewFuncCollection creates a new FuncCollection with the given functions.
func NewFuncCollection(funcs ...*Func) FuncCollection {
	return FuncCollection{}.Register(funcs...)
}

// Register registers the given functions with the FuncCollection.
func (fc FuncCollection) Register(funcs ...*Func) FuncCollection {
	for _, f := range funcs {
		fc[f.name] = f.Handler()
	}
	return fc
}

// Get returns the function with the given name.
func (fc FuncCollection) Get(name string) (FuncFn, bool) {
	fn, ok := fc[name]
	return fn, ok
}

// Delete deletes the functions with the given names.
func (fc FuncCollection) Delete(names ...string) FuncCollection {
	for _, name := range names {
		delete(fc, name)
	}
	return fc
}

// Copy returns a copy of the FuncCollection.
func (fc FuncCollection) Copy() FuncCollection {
	c := NewFuncCollection()
	for k, v := range fc {
		c[k] = v
	}
	return c
}
