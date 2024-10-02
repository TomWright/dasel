package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
)

func spreadExprExecutor() (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		s := model.NewSliceValue()

		s.MarkAsSpread()

		switch {
		case data.IsSlice():
			v, err := data.SliceValue()
			if err != nil {
				return nil, fmt.Errorf("error getting slice value: %w", err)
			}
			for _, sv := range v {
				if err := s.Append(model.NewValue(sv)); err != nil {
					return nil, fmt.Errorf("error appending value to slice: %w", err)
				}
			}
		case data.IsMap():
			if err := data.RangeMap(func(key string, value *model.Value) error {
				if err := s.Append(value); err != nil {
					return fmt.Errorf("error appending value to slice: %w", err)
				}
				return nil
			}); err != nil {
				return nil, fmt.Errorf("error ranging map: %w", err)
			}
		default:
			return nil, fmt.Errorf("cannot spread on type %s", data.Type())
		}

		return s, nil
	}, nil
}

// prepareSpreadValues looks at the incoming value, and if we detect a spread value, we return the individual values.
func prepareSpreadValues(val *model.Value) (model.Values, error) {
	if val.IsSlice() && val.IsSpread() {
		sliceLen, err := val.SliceLen()
		if err != nil {
			return nil, fmt.Errorf("error getting slice length: %w", err)
		}
		values := make(model.Values, sliceLen)
		for i := 0; i < sliceLen; i++ {
			v, err := val.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index %d: %w", i, err)
			}
			values[i] = v
		}
		return values, nil
	}
	return model.Values{val}, nil
}
