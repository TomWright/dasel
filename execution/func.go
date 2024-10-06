package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

type FuncFn func(data *model.Value, args model.Values) (*model.Value, error)

var singleResponseFuncLookup = map[string]FuncFn{}

func RegisterFunc(name string, fn FuncFn) {
	singleResponseFuncLookup[name] = fn
}

func init() {
	RegisterFunc("len", func(data *model.Value, args model.Values) (*model.Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("len expects a single argument")
		}

		arg := args[0]

		l, err := arg.Len()
		if err != nil {
			return nil, err
		}

		return model.NewIntValue(int64(l)), nil
	})

	RegisterFunc("add", func(_ *model.Value, args model.Values) (*model.Value, error) {
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
	})
}
