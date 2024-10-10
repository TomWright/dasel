package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func conditionalExprExecutor(e ast.ConditionalExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		cond, err := ExecuteAST(e.Cond, data)
		if err != nil {
			return nil, fmt.Errorf("error evaluating condition: %w", err)
		}

		condBool, err := cond.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error converting condition to boolean: %w", err)
		}

		if condBool {
			res, err := ExecuteAST(e.Then, data)
			if err != nil {
				return nil, fmt.Errorf("error executing then block: %w", err)
			}
			return res, nil
		}

		if e.Else != nil {
			res, err := ExecuteAST(e.Else, data)
			if err != nil {
				return nil, fmt.Errorf("error executing else block: %w", err)
			}
			return res, nil
		}

		return model.NewNullValue(), nil
	}, nil
}
