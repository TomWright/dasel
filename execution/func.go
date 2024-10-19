package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// ArgsValidator is a function that validates the arguments passed to a function.
type ArgsValidator func(name string, args model.Values) error

// ValidateArgsExactly returns an ArgsValidator that validates that the number of arguments passed to a function is exactly the expected number.
func ValidateArgsExactly(expected int) ArgsValidator {
	return func(name string, args model.Values) error {
		if len(args) == expected {
			return nil
		}
		return fmt.Errorf("func %q expects exactly %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMin returns an ArgsValidator that validates that the number of arguments passed to a function is at least the expected number.
func ValidateArgsMin(expected int) ArgsValidator {
	return func(name string, args model.Values) error {
		if len(args) >= expected {
			return nil
		}
		return fmt.Errorf("func %q expects at least %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMax returns an ArgsValidator that validates that the number of arguments passed to a function is at most the expected number.
func ValidateArgsMax(expected int) ArgsValidator {
	return func(name string, args model.Values) error {
		if len(args) <= expected {
			return nil
		}
		return fmt.Errorf("func %q expects no more than %d arguments, got %d", name, expected, len(args))
	}
}

// ValidateArgsMinMax returns an ArgsValidator that validates that the number of arguments passed to a function is between the min and max expected numbers.
func ValidateArgsMinMax(min int, max int) ArgsValidator {
	return func(name string, args model.Values) error {
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
	return func(data *model.Value, args model.Values) (*model.Value, error) {
		if f.argsValidator != nil {
			if err := f.argsValidator(f.name, args); err != nil {
				return nil, err
			}
		}
		return f.handler(data, args)
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
type FuncFn func(data *model.Value, args model.Values) (*model.Value, error)

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

var (
	// DefaultFuncCollection is the default collection of functions that can be executed.
	DefaultFuncCollection = NewFuncCollection(
		FuncLen,
		FuncAdd,
		FuncToString,
		FuncMerge,
	)

	// FuncLen is a function that returns the length of the given value.
	FuncLen = NewFunc(
		"len",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			arg := args[0]

			l, err := arg.Len()
			if err != nil {
				return nil, err
			}

			return model.NewIntValue(int64(l)), nil
		},
		ValidateArgsExactly(1),
	)

	// FuncAdd is a function that adds the given values together.
	FuncAdd = NewFunc(
		"add",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			var foundInts, foundFloats int
			var intRes int64
			var floatRes float64
			for _, arg := range args {
				if arg.IsFloat() {
					foundFloats++
					v, err := arg.FloatValue()
					if err != nil {
						return nil, fmt.Errorf("error getting float value: %w", err)
					}
					floatRes += v
					continue
				}
				if arg.IsInt() {
					foundInts++
					v, err := arg.IntValue()
					if err != nil {
						return nil, fmt.Errorf("error getting int value: %w", err)
					}
					intRes += v
					continue
				}
				return nil, fmt.Errorf("expected int or float, got %s", arg.Type())
			}
			if foundFloats > 0 {
				return model.NewFloatValue(floatRes + float64(intRes)), nil
			}
			return model.NewIntValue(intRes), nil
		},
		ValidateArgsMin(1),
	)

	// FuncToString is a function that converts the given value to a string.
	FuncToString = NewFunc(
		"toString",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			switch args[0].Type() {
			case model.TypeString:
				return args[0], nil
			case model.TypeInt:
				i, err := args[0].IntValue()
				if err != nil {
					return nil, err
				}
				return model.NewStringValue(fmt.Sprintf("%d", i)), nil
			case model.TypeFloat:
				i, err := args[0].FloatValue()
				if err != nil {
					return nil, err
				}
				return model.NewStringValue(fmt.Sprintf("%f", i)), nil
			case model.TypeBool:
				i, err := args[0].BoolValue()
				if err != nil {
					return nil, err
				}
				return model.NewStringValue(fmt.Sprintf("%v", i)), nil
			default:
				return nil, fmt.Errorf("cannot convert %s to string", args[0].Type())
			}
		},
		ValidateArgsExactly(1),
	)

	// FuncMerge is a function that merges two or more items together.
	FuncMerge = NewFunc(
		"merge",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			if len(args) == 1 {
				return args[0], nil
			}

			expectedType := args[0].Type()

			switch expectedType {
			case model.TypeMap:
				break
			default:
				return nil, fmt.Errorf("merge exects a map, found %s", expectedType)
			}

			// Validate types match
			for _, a := range args {
				if a.Type() != expectedType {
					return nil, fmt.Errorf("merge expects all arguments to be of the same type. expected %s, got %s", expectedType.String(), a.Type().String())
				}
			}

			base := model.NewMapValue()

			for i := 0; i < len(args); i++ {
				next := args[i]

				nextKVs, err := next.MapKeyValues()
				if err != nil {
					return nil, fmt.Errorf("merge failed to extract key values for arg %d: %w", i, err)
				}

				for _, kv := range nextKVs {
					if err := base.SetMapKey(kv.Key, kv.Value); err != nil {
						return nil, fmt.Errorf("merge failed to set map key %s: %w", kv.Key, err)
					}
				}
			}

			return base, nil
		},
		ValidateArgsMin(1),
	)
)
