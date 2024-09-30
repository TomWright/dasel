package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func mapExprExecutor(e ast.MapExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		if !data.IsSlice() {
			return nil, fmt.Errorf("cannot map over non-array")
		}
		sliceLen, err := data.SliceLen()
		if err != nil {
			return nil, fmt.Errorf("error getting slice length: %w", err)
		}
		res := model.NewSliceValue()

		for i := 0; i < sliceLen; i++ {
			item, err := data.GetSliceIndex(i)
			if err != nil {
				return nil, fmt.Errorf("error getting slice index: %w", err)
			}
			for _, expr := range e.Exprs {
				item, err = ExecuteAST(expr, item)
				if err != nil {
					return nil, err
				}
			}
			res.Append(item)
		}

		return res, nil
	}, nil
}
