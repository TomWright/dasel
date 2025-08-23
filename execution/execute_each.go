package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func eachExprExecutor(e ast.EachExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot each over non-array")
		}

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			_, err := ExecuteAST(e.Expr, item, options)
			if err != nil {
				return err
			}
			// Each disregards the output.
			return nil
		}); err != nil {
			return nil, fmt.Errorf("error ranging over slice: %w", err)
		}

		return data, nil
	}, nil
}
