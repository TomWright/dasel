package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

func binaryExprExecutor(opts *Options, e ast.BinaryExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		left, err := ExecuteAST(e.Left, data, opts)
		if err != nil {
			return nil, fmt.Errorf("error evaluating left expression: %w", err)
		}
		right, err := ExecuteAST(e.Right, data, opts)
		if err != nil {
			return nil, fmt.Errorf("error evaluating right expression: %w", err)
		}

		switch e.Operator.Kind {
		case lexer.Plus:
			return left.Add(right)
		case lexer.Dash:
			return left.Subtract(right)
		case lexer.Star:
			return left.Multiply(right)
		case lexer.Slash:
			return left.Divide(right)
		case lexer.Percent:
			return left.Modulo(right)
		case lexer.GreaterThan:
			return left.GreaterThan(right)
		case lexer.GreaterThanOrEqual:
			return left.GreaterThanOrEqual(right)
		case lexer.LessThan:
			return left.LessThan(right)
		case lexer.LessThanOrEqual:
			return left.LessThanOrEqual(right)
		case lexer.Equal:
			return left.Equal(right)
		case lexer.NotEqual:
			return left.NotEqual(right)
		case lexer.Equals:
			err := left.Set(right)
			return left, err
		case lexer.And:
			leftBool, err := left.BoolValue()
			if err != nil {
				return nil, fmt.Errorf("error getting left bool value: %w", err)
			}
			rightBool, err := right.BoolValue()
			if err != nil {
				return nil, fmt.Errorf("error getting right bool value: %w", err)
			}
			return model.NewBoolValue(leftBool && rightBool), nil
		case lexer.Or:
			leftBool, err := left.BoolValue()
			if err != nil {
				return nil, fmt.Errorf("error getting left bool value: %w", err)
			}
			rightBool, err := right.BoolValue()
			if err != nil {
				return nil, fmt.Errorf("error getting right bool value: %w", err)
			}
			return model.NewBoolValue(leftBool || rightBool), nil
		case lexer.Like, lexer.NotLike:
			leftStr, err := left.StringValue()
			if err != nil {
				return nil, fmt.Errorf("like requires left side to be a string, got %s", left.Type().String())
			}
			rightPatt, ok := e.Right.(ast.RegexExpr)
			if !ok {
				return nil, fmt.Errorf("like requires right side to be a regex pattern")
			}
			res := rightPatt.Regex.MatchString(leftStr)
			if e.Operator.Kind == lexer.NotLike {
				res = !res
			}
			return model.NewBoolValue(res), nil
		default:
			return nil, fmt.Errorf("unhandled operator: %s", e.Operator.Value)
		}
	}, nil
}
