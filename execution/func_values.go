package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncValues returns the values of a map as an array.
var FuncValues = NewFunc(
	"values",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		if !data.IsMap() {
			return nil, fmt.Errorf("values can only be used on maps, got %s", data.Type().String())
		}

		res := model.NewSliceValue()
		if err := data.RangeMap(func(key string, value *model.Value) error {
			return res.Append(value)
		}); err != nil {
			return nil, err
		}

		return res, nil
	},
	ValidateArgsExactly(0),
)
