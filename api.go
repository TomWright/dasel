// Package dasel contains everything you'll need to use dasel from a go application.
package dasel

import (
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

// Query queries the data using the selector and returns the results.
func Query(data any, selector string, opts ...execution.ExecuteOptionFn) ([]*model.Value, int, error) {
	options := execution.NewOptions(opts...)
	val := model.NewValue(data)
	out, err := execution.ExecuteSelector(selector, val, options)
	if err != nil {
		return nil, 0, err
	}

	if out.IsBranch() {
		res := make([]*model.Value, 0)
		if err := out.RangeSlice(func(i int, v *model.Value) error {
			res = append(res, v)
			return nil
		}); err != nil {
			return nil, 0, err
		}
		return res, len(res), nil
	}

	return []*model.Value{out}, 1, nil
}

func Select(data any, selector string, opts ...execution.ExecuteOptionFn) (any, int, error) {
	res, count, err := Query(data, selector, opts...)
	if err != nil {
		return nil, 0, err
	}
	out := make([]any, 0)
	for _, v := range res {
		out = append(out, v.Interface())
	}
	return out, count, err
}

func Modify(data any, selector string, newValue any, opts ...execution.ExecuteOptionFn) (int, error) {
	res, count, err := Query(data, selector, opts...)
	if err != nil {
		return 0, err
	}
	for _, v := range res {
		if err := v.Set(model.NewValue(newValue)); err != nil {
			return 0, err
		}
	}
	return count, nil
}
