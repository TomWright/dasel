package execution

import (
	"context"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func searchExprExecutor(e ast.SearchExpr) (expressionExecutor, error) {
	var doSearch func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error)
	doSearch = func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error) {
		res := make([]*model.Value, 0)

		switch data.Type() {
		case model.TypeMap:
			if err := data.RangeMap(func(key string, v *model.Value) error {
				got, err := ExecuteAST(ctx, e.Expr, v, options)
				if err != nil {
					return err
				}

				gotV, err := got.BoolValue()
				if err != nil {
					return err
				}

				if gotV {
					res = append(res, v)
				}

				gotNext, err := doSearch(ctx, options, v)
				if err != nil {
					return err
				}
				res = append(res, gotNext...)

				return nil
			}); err != nil {
				return nil, err
			}
		case model.TypeSlice:
			if err := data.RangeSlice(func(i int, v *model.Value) error {
				got, err := ExecuteAST(ctx, e.Expr, v, options)
				if err != nil {
					return err
				}

				gotV, err := got.BoolValue()
				if err != nil {
					return err
				}

				if gotV {
					res = append(res, v)
				}

				gotNext, err := doSearch(ctx, options, v)
				if err != nil {
					return err
				}
				res = append(res, gotNext...)

				return nil
			}); err != nil {
				return nil, err
			}
		}

		return res, nil
	}

	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "searchExpr")
		matches := model.NewSliceValue()

		found, err := doSearch(ctx, options, data)
		if err != nil {
			return nil, err
		}

		for _, f := range found {
			if err := matches.Append(f); err != nil {
				return nil, err
			}
		}

		return matches, nil
	}, nil
}
