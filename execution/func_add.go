package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncAdd is a function that adds the given values together.
var FuncAdd = NewFunc(
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
