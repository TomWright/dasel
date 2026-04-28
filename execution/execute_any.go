package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func anyExprExecutor(e ast.AnyExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "anyExpr")
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot use any over non-array")
		}

		result := false
		if err := data.RangeSlice(func(i int, item *model.Value) error {
			v, err := ExecuteAST(ctx, e.Expr, item, options)
			if err != nil {
				return err
			}

			boolV, err := v.BoolValue()
			if err != nil {
				return err
			}

			if boolV {
				result = true
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		return model.NewBoolValue(result), nil
	}, nil
}
