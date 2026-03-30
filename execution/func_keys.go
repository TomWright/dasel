package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

// FuncKeys returns the keys of a map or the indices of a slice.
var FuncKeys = NewFunc(
	"keys",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		switch data.Type() {
		case model.TypeMap:
			keys, err := data.MapKeys()
			if err != nil {
				return nil, err
			}
			res := model.NewSliceValue()
			for _, key := range keys {
				if err := res.Append(model.NewStringValue(key)); err != nil {
					return nil, err
				}
			}
			return res, nil
		case model.TypeSlice:
			len, err := data.SliceLen()
			if err != nil {
				return nil, err
			}
			res := model.NewSliceValue()
			for i := 0; i < len; i++ {
				if err := res.Append(model.NewIntValue(int64(i))); err != nil {
					return nil, err
				}
			}
			return res, nil
		default:
			return nil, fmt.Errorf("keys can only be used on maps and slices")
		}
	},
	ValidateArgsExactly(0),
)
