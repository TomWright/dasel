package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func allExprExecutor(e ast.AllExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "allExpr")
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot use all over non-array")
		}

		result := true
		if err := data.RangeSlice(func(i int, item *model.Value) error {
			v, err := ExecuteAST(ctx, e.Expr, item, options)
			if err != nil {
				return err
			}

			boolV, err := v.BoolValue()
			if err != nil {
				return err
			}

			if !boolV {
				result = false
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		return model.NewBoolValue(result), nil
	}, nil
}
