package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func branchExprExecutor(opts *Options, e ast.BranchExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		res := model.NewSliceValue()

		for _, expr := range e.Exprs {
			r, err := ExecuteAST(expr, data, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to execute branch expr: %w", err)
			}

			// This deals with the spread operator in the branch expression.
			valsToAppend, err := prepareSpreadValues(r)
			if err != nil {
				return nil, fmt.Errorf("error handling spread values: %w", err)
			}
			for _, v := range valsToAppend {
				if err := res.Append(v); err != nil {
					return nil, fmt.Errorf("failed to append branch result: %w", err)
				}
			}
		}

		res.MarkAsBranch()

		return res, nil
	}, nil
}
