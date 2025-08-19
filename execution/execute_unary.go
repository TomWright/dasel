package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func unaryExprExecutor(opts *Options, e ast.UnaryExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		right, err := ExecuteAST(e.Right, data, opts)
		if err != nil {
			return nil, fmt.Errorf("error evaluating right expression: %w", err)
		}

		switch e.Operator.Kind {
		case lexer.Exclamation:
			boolV, err := right.BoolValue()
			if err != nil {
				return nil, fmt.Errorf("error converting value to boolean: %w", err)
			}
			return model.NewBoolValue(!boolV), nil
		default:
			return nil, fmt.Errorf("unhandled unary operator: %s", e.Operator.Value)
		}
	}, nil
}
