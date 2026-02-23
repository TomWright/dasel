package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func arrayExprExecutor(e ast.ArrayExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "arrayExpr")
		res := model.NewSliceValue()

		for _, expr := range e.Exprs {
			el, err := ExecuteAST(ctx, expr, data, options)
			if err != nil {
				return nil, err
			}

			if el.IsSpread() {
				if err := el.RangeSlice(func(_ int, value *model.Value) error {
					if err := res.Append(value); err != nil {
						return err
					}
					return nil
				}); err != nil {
					return nil, err
				}
			} else {
				if err := res.Append(el); err != nil {
					return nil, err
				}
			}
		}

		return res, nil
	}, nil
}

func rangeExprExecutor(e ast.RangeExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "rangeExpr")
		var start, end int64 = 0, -1
		if e.Start != nil {
			startE, err := ExecuteAST(ctx, e.Start, data, options)
			if err != nil {
				return nil, fmt.Errorf("error evaluating start expression: %w", err)
			}

			start, err = startE.IntValue()
			if err != nil {
				return nil, fmt.Errorf("error getting start int value: %w", err)
			}
		}

		if e.End != nil {
			endE, err := ExecuteAST(ctx, e.End, data, options)
			if err != nil {
				return nil, fmt.Errorf("error evaluating end expression: %w", err)
			}

			end, err = endE.IntValue()
			if err != nil {
				return nil, fmt.Errorf("error getting end int value: %w", err)
			}
		}

		var res *model.Value
		var err error

		switch data.Type() {
		case model.TypeString:
			res, err = data.StringIndexRange(int(start), int(end))
		case model.TypeSlice:
			res, err = data.SliceIndexRange(int(start), int(end))
		default:
			err = fmt.Errorf("range expects a slice or string, got %s", data.Type())
		}

		if err != nil {
			return nil, err
		}

		return res, nil
	}, nil
}

func indexExprExecutor(e ast.IndexExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "indexExpr")
		indexE, err := ExecuteAST(ctx, e.Index, data, options)
		if err != nil {
			return nil, fmt.Errorf("error evaluating index expression: %w", err)
		}

		index, err := indexE.IntValue()
		if err != nil {
			return nil, fmt.Errorf("error getting index int value: %w", err)
		}

		return data.GetSliceIndex(int(index))
	}, nil
}
