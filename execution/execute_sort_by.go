package execution

import (
	"context"
	"fmt"
	"slices"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func sortByExprExecutor(e ast.SortByExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "sortByExpr")
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot sort by on non-slice data")
		}

		type sortableValue struct {
			index int
			value *model.Value
		}
		values := make([]sortableValue, 0)

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			item, err := ExecuteAST(ctx, e.Expr, item, options)
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
