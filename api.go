// Package dasel contains everything you'll need to use dasel from a go application.
package dasel

import (
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func Select(data any, selector string) (any, error) {
	val := model.NewValue(data)
	res, err := execution.ExecuteSelector(selector, val)
	if err != nil {
		return nil, err
	}
	return res.Interface(), nil
}

func Modify(data any, selector string, newValue any) error {
	val := model.NewValue(data)
	newVal := model.NewValue(newValue)
	res, err := execution.ExecuteSelector(selector, val)
	if err != nil {
		return err
	}

	if err := res.Set(newVal); err != nil {
		return err
	}

	return nil
}
