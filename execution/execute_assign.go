package execution

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func executeAssign(left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
	err := left.Set(right)
	if err != nil {
		return nil, fmt.Errorf("error setting value: %w", err)
	}
	return right, nil
}

func variableAssignExprExecutor(e ast.AssignExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		varName := e.Variable.Name

		value, err := ExecuteAST(e.Value, data, options)
		if err != nil {
			return nil, fmt.Errorf("error executing variable assignment expression: %w", err)
		}

		//if varName == "this" {
		//	return executeAssign(data, value, ast.BinaryExpr{})
		//}

		_, ok := options.Vars[varName]
		if !ok {
			options.Vars[varName] = value
		} else {
			if err := data.Set(value); err != nil {
				return nil, fmt.Errorf("error updating variable assignment expression: %w", err)
			}
		}

		return value, nil
	}, nil
}
