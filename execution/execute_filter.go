package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func filterExprExecutor(opts *Options, e ast.FilterExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot filter over non-array")
		}
		res := model.NewSliceValue()

		if err := data.RangeSlice(func(i int, item *model.Value) error {
			v, err := ExecuteAST(e.Expr, item, opts)
			if err != nil {
				return err
			}

			boolV, err := v.BoolValue()
			if err != nil {
				return err
			}

			if !boolV {
				return nil
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
