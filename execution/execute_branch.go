package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func branchExprExecutor(e ast.BranchExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		res := model.NewSliceValue()
		res.MarkAsBranch()

		if len(e.Exprs) == 0 {
			// No expressions given. We'll branch on the input data.
			if err := data.RangeSlice(func(_ int, value *model.Value) error {
				if err := res.Append(value); err != nil {
					return fmt.Errorf("failed to append branch result: %w", err)
				}
				return nil
			}); err != nil {
				return nil, fmt.Errorf("failed to range slice: %w", err)
			}
		} else {
			for _, expr := range e.Exprs {
				r, err := ExecuteAST(expr, data, options)
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
		}

		return res, nil
	}, nil
}
