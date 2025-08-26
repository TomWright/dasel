package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func mapExprExecutor(e ast.MapExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "mapExpr")
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot map over non-array")
		}
		res := model.NewSliceValue()

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			item, err := ExecuteAST(ctx, e.Expr, item, options)
			if err != nil {
				return err
			}
			if err := res.Append(item); err != nil {
				return fmt.Errorf("error appending item to result: %w", err)
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}
		return res, nil
	}, nil
}
