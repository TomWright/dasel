package execution

import (
	"fmt"
	"slices"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func sortByExprExecutor(e ast.SortByExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot sort by on non-slice data")
		}

		type sortableValue struct {
			index int
			value *model.Value
		}
		values := make([]sortableValue, 0)

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			item, err := ExecuteAST(e.Expr, item, options)
			if err != nil {
				return err
			}
			values = append(values, sortableValue{
				index: i,
				value: item,
			})
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		slices.SortFunc(values, func(i, j sortableValue) int {
			res, err := i.value.Compare(j.value)
			if err != nil {
				return 0
			}
			if e.Descending {
				return -res
			}
			return res
		})

		res := model.NewSliceValue()

		for _, i := range values {
			item, err := data.GetSliceIndex(i.index)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			if err := res.Append(item); err != nil {
				return nil, fmt.Errorf("error appending item to result: %w", err)
			}
		}

		return res, nil
	}, nil
}

func sortByExprExecutor2(e ast.SortByExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot sort by on non-slice data")
		}

		sortedValues := model.NewSliceValue()
		sortedIndexes := make([]int, 0)

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			item, err := ExecuteAST(e.Expr, item, options)
			if err != nil {
				return err
			}
			if err := sortedValues.Append(item); err != nil {
				return fmt.Errorf("error appending item to result: %w", err)
			}
			sortedIndexes = append(sortedIndexes, i)
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		l, err := sortedValues.Len()
		if err != nil {
			return nil, fmt.Errorf("error getting length of slice: %w", err)
		}

		for i := 0; i < l-1; i++ {
			cur, err := sortedValues.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			curIndex := sortedIndexes[i]
			next, err := sortedValues.GetSliceIndex(i + 1)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			nextIndex := sortedIndexes[i+1]

			cmp, err := cur.Compare(next)
			if err != nil {
				return nil, fmt.Errorf("error comparing values: %w", err)
			}

			if cmp == 0 {
				continue
			}

			if !e.Descending {
				if cmp > 0 {
					if err := sortedValues.SetSliceIndex(i, next); err != nil {
						return nil, fmt.Errorf("error setting slice index: %w", err)
					}
					sortedIndexes[i] = nextIndex
					if err := sortedValues.SetSliceIndex(i+1, cur); err != nil {
						return nil, fmt.Errorf("error setting slice index: %w", err)
					}
					sortedIndexes[i+1] = curIndex
					i -= 1
				}
			} else {
				if cmp < 0 {
					if err := sortedValues.SetSliceIndex(i, next); err != nil {
						return nil, fmt.Errorf("error setting slice index: %w", err)
					}
					sortedIndexes[i] = nextIndex
					if err := sortedValues.SetSliceIndex(i+1, cur); err != nil {
						return nil, fmt.Errorf("error setting slice index: %w", err)
					}
					sortedIndexes[i+1] = curIndex
					i -= 1
				}
			}
		}

		res := model.NewSliceValue()

		for _, i := range sortedIndexes {
			item, err := data.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			if err := res.Append(item); err != nil {
				return nil, fmt.Errorf("error appending item to result: %w", err)
			}
		}

		return res, nil
	}, nil
}
