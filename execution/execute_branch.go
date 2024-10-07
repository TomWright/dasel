package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func branchExprExecutor(e ast.BranchExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		res := model.NewSliceValue()

		for _, expr := range e.Exprs {
			r, err := ExecuteAST(expr, data)
			if err != nil {
				return nil, fmt.Errorf("failed to execute branch expr: %w", err)
			}
			if err := res.Append(r); err != nil {
				return nil, fmt.Errorf("failed to append branch result: %w", err)
			}
		}

		res.MarkAsBranch()

		return res, nil
	}, nil
}
