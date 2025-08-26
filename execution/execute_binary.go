package execution

import (
	"context"
	"errors"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
)

type binaryExpressionExecutorFn func(ctx context.Context, expr ast.BinaryExpr, value *model.Value, options *Options) (*model.Value, error)

func basicBinaryExpressionExecutorFn(handler func(ctx context.Context, left *model.Value, right *model.Value, e ast.BinaryExpr) (*model.Value, error)) binaryExpressionExecutorFn {
	return func(ctx context.Context, expr ast.BinaryExpr, value *model.Value, options *Options) (*model.Value, error) {
		left, err := ExecuteAST(ctx, expr.Left, value, options)
		if err != nil {
			return nil, fmt.Errorf("error evaluating left expression: %w", err)
		}

		if !left.IsBranch() {
			right, err := ExecuteAST(ctx, expr.Right, value, options)
			if err != nil {
				return nil, fmt.Errorf("error evaluating right expression: %w", err)
			}
			res, err := handler(ctx, left, right, expr)
			if err != nil {
				return nil, err
			}
			return res, nil
		}

		res := model.NewSliceValue()
		res.MarkAsBranch()
		if err := left.RangeSlice(func(i int, v *model.Value) error {
			right, err := ExecuteAST(ctx, expr.Right, v, options)
			if err != nil {
				return fmt.Errorf("error evaluating right expression: %w", err)
			}
			r, err := handler(ctx, v, right, expr)
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
	}
}

var binaryExpressionExecutors = map[lexer.TokenKind]binaryExpressionExecutorFn{}

func binaryExprExecutor(e ast.BinaryExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "binaryExpr")
		if e.Left == nil || e.Right == nil {
			return nil, fmt.Errorf("left and right expressions must be provided")
		}

		exec, ok := binaryExpressionExecutors[e.Operator.Kind]
		if !ok {
			return nil, fmt.Errorf("unhandled operator: %s", e.Operator.Value)
		}

		return exec(ctx, e, data, options)
	}, nil
}

func init() {
	binaryExpressionExecutors[lexer.Plus] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Add(right)
	})
	binaryExpressionExecutors[lexer.Dash] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Subtract(right)
	})
	binaryExpressionExecutors[lexer.Star] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Multiply(right)
	})
	binaryExpressionExecutors[lexer.Slash] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Divide(right)
	})
	binaryExpressionExecutors[lexer.Percent] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Modulo(right)
	})
	binaryExpressionExecutors[lexer.GreaterThan] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.GreaterThan(right)
	})
	binaryExpressionExecutors[lexer.GreaterThanOrEqual] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.GreaterThanOrEqual(right)
	})
	binaryExpressionExecutors[lexer.LessThan] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.LessThan(right)
	})
	binaryExpressionExecutors[lexer.LessThanOrEqual] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.LessThanOrEqual(right)
	})
	binaryExpressionExecutors[lexer.Equal] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.Equal(right)
	})
	binaryExpressionExecutors[lexer.NotEqual] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		return left.NotEqual(right)
	})
	binaryExpressionExecutors[lexer.Equals] = func(ctx context.Context, expr ast.BinaryExpr, value *model.Value, options *Options) (*model.Value, error) {
		if leftVar, ok := expr.Left.(ast.VariableExpr); ok {
			// It is expected that the left side of an assignment may not exist yet.
			if _, ok := options.Vars[leftVar.Name]; !ok {
				options.Vars[leftVar.Name] = model.NewNullValue()
			}
		}
		return basicBinaryExpressionExecutorFn(executeAssign)(ctx, expr, value, options)
	}
	binaryExpressionExecutors[lexer.And] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		leftBool, err := left.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error getting left bool value: %w", err)
		}
		rightBool, err := right.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error getting right bool value: %w", err)
		}
		return model.NewBoolValue(leftBool && rightBool), nil
	})
	binaryExpressionExecutors[lexer.Or] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
		leftBool, err := left.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error getting left bool value: %w", err)
		}
		rightBool, err := right.BoolValue()
		if err != nil {
			return nil, fmt.Errorf("error getting right bool value: %w", err)
		}
		return model.NewBoolValue(leftBool || rightBool), nil
	})
	binaryExpressionExecutors[lexer.Like] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, _ *model.Value, e ast.BinaryExpr) (*model.Value, error) {
		leftStr, err := left.StringValue()
		if err != nil {
			return nil, fmt.Errorf("like requires left side to be a string, got %s", left.Type().String())
		}
		rightPatt, ok := e.Right.(ast.RegexExpr)
		if !ok {
			return nil, fmt.Errorf("like requires right side to be a regex pattern")
		}
		res := rightPatt.Regex.MatchString(leftStr)
		return model.NewBoolValue(res), nil
	})
	binaryExpressionExecutors[lexer.NotLike] = basicBinaryExpressionExecutorFn(func(ctx context.Context, left *model.Value, _ *model.Value, e ast.BinaryExpr) (*model.Value, error) {
		leftStr, err := left.StringValue()
		if err != nil {
			return nil, fmt.Errorf("like requires left side to be a string, got %s", left.Type().String())
		}
		rightPatt, ok := e.Right.(ast.RegexExpr)
		if !ok {
			return nil, fmt.Errorf("like requires right side to be a regex pattern")
		}
		res := rightPatt.Regex.MatchString(leftStr)
		return model.NewBoolValue(!res), nil
	})
	binaryExpressionExecutors[lexer.DoubleQuestionMark] = func(ctx context.Context, expr ast.BinaryExpr, value *model.Value, options *Options) (*model.Value, error) {
		left, err := ExecuteAST(ctx, expr.Left, value, options)

		if err == nil && !left.IsNull() {
			return left, nil
		}

		if err != nil {
			handleErrs := []any{
				model.ErrIncompatibleTypes{},
				model.ErrUnexpectedType{},
				model.ErrUnexpectedTypes{},
				model.SliceIndexOutOfRange{},
				model.MapKeyNotFound{},
			}
			for _, e := range handleErrs {
				if errors.As(err, &e) {
					err = nil
					break
				}
			}

			if err != nil {
				return nil, fmt.Errorf("error evaluating left expression: %w", err)
			}
		}

		// Do we need to handle branches here?
		right, err := ExecuteAST(ctx, expr.Right, value, options)
		if err != nil {
			return nil, fmt.Errorf("error evaluating right expression: %w", err)
		}
		return right, nil
	}
}
