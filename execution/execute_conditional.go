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
			return ExecuteAST(e.Then, data)
		}

		if e.Else != nil {
			return ExecuteAST(e.Else, data)
		}

		return model.NewNullValue(), nil
	}, nil
}
