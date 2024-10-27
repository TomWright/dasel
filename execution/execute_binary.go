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

		var doOperation func(a *model.Value, b *model.Value) (*model.Value, error)

		switch e.Operator.Kind {
		case lexer.Plus:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Add(b)
			}
		case lexer.Dash:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Subtract(b)
			}
		case lexer.Star:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Multiply(b)
			}
		case lexer.Slash:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Divide(b)
			}
		case lexer.Percent:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Modulo(b)
			}
		case lexer.GreaterThan:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.GreaterThan(b)
			}
		case lexer.GreaterThanOrEqual:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.GreaterThanOrEqual(b)
			}
		case lexer.LessThan:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.LessThan(b)
			}
		case lexer.LessThanOrEqual:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.LessThanOrEqual(b)
			}
		case lexer.Equal:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.Equal(b)
			}
		case lexer.NotEqual:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				return a.NotEqual(b)
			}
		case lexer.Equals:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				err := a.Set(b)
				if err != nil {
					return nil, fmt.Errorf("error setting value: %w", err)
				}
				switch a.Type() {
				case model.TypeMap:
					return a, nil
				case model.TypeSlice:
					return a, nil
				default:
					return b, nil
				}
			}
		case lexer.And:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				leftBool, err := a.BoolValue()
				if err != nil {
					return nil, fmt.Errorf("error getting left bool value: %w", err)
				}
				rightBool, err := b.BoolValue()
				if err != nil {
					return nil, fmt.Errorf("error getting right bool value: %w", err)
				}
				return model.NewBoolValue(leftBool && rightBool), nil
			}
		case lexer.Or:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				leftBool, err := a.BoolValue()
				if err != nil {
					return nil, fmt.Errorf("error getting left bool value: %w", err)
				}
				rightBool, err := b.BoolValue()
				if err != nil {
					return nil, fmt.Errorf("error getting right bool value: %w", err)
				}
				return model.NewBoolValue(leftBool || rightBool), nil
			}
		case lexer.Like, lexer.NotLike:
			doOperation = func(a *model.Value, b *model.Value) (*model.Value, error) {
				leftStr, err := a.StringValue()
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
			}
		default:
			return nil, fmt.Errorf("unhandled operator: %s", e.Operator.Value)
		}

		if doOperation == nil {
			return nil, fmt.Errorf("missing operation for operator %s", e.Operator.Value)
		}

		if !left.IsBranch() {
			right, err := ExecuteAST(e.Right, data, opts)
			if err != nil {
				return nil, fmt.Errorf("error evaluating right expression: %w", err)
			}
			return doOperation(left, right)
		}

		res := model.NewSliceValue()
		res.MarkAsBranch()
		if err := left.RangeSlice(func(i int, v *model.Value) error {
			right, err := ExecuteAST(e.Right, v, opts)
			if err != nil {
				return fmt.Errorf("error evaluating right expression: %w", err)
			}

			r, err := doOperation(v, right)
			if err != nil {
				return err
			}
			if err := res.Append(r); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return res, nil
	}, nil
}
