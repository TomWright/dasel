package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func conditionalExprExecutor(e ast.ConditionalExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "conditionalExpr")
		cond, err := ExecuteAST(ctx, e.Cond, data, options)
		if err != nil {
			return nil, fmt.Errorf("error evaluating condition: %w", err)
		}

		condBool, err := cond.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error converting condition to boolean: %w", err)
		}

		if condBool {
			res, err := ExecuteAST(ctx, e.Then, data, options)
			if err != nil {
				return nil, fmt.Errorf("error executing then block: %w", err)
			}
			return res, nil
		}

		if e.Else != nil {
			res, err := ExecuteAST(ctx, e.Else, data, options)
			if err != nil {
				return nil, fmt.Errorf("error executing else block: %w", err)
			}
			return res, nil
		}

		return model.NewNullValue(), nil
	}, nil
}
