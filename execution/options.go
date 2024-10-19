package execution

import "github.com/tomwright/dasel/v3/model"

// ExecuteOptionFn is a function that can be used to set options on the execution of the selector.
type ExecuteOptionFn func(*Options)

// Options contains the options for the execution of the selector.
type Options struct {
	Funcs FuncCollection
	Vars  map[string]*model.Value
}

// NewOptions creates a new Options struct with the given options.
func NewOptions(opts ...ExecuteOptionFn) *Options {
	o := &Options{
		Funcs: DefaultFuncCollection,
		Vars:  map[string]*model.Value{},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(o)
	}
	return o
}

// WithFuncs sets the functions that can be used in the selector.
func WithFuncs(fc FuncCollection) ExecuteOptionFn {
	return func(o *Options) {
		o.Funcs = fc
	}
}

// WithVariable sets a variable for use in the selector.
func WithVariable(key string, val *model.Value) ExecuteOptionFn {
	return func(o *Options) {
		o.Vars[key] = val
	}
}
