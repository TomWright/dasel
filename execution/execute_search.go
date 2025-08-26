package execution

import (
	"context"
	"errors"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func searchExprExecutor(e ast.SearchExpr) (expressionExecutor, error) {
	var doSearch func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error)
	processValue := func(ctx context.Context, v *model.Value, options *Options) (bool, error) {
		got, err := ExecuteAST(ctx, e.Expr, v, options)
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
		}
		if err != nil {
			return false, err
		}

		if got == nil {
			return false, nil
		}

		gotV, err := got.BoolValue()
		if err != nil {
			return false, err
		}
		return gotV, nil
	}
	doSearch = func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error) {
		res := make([]*model.Value, 0)

		switch data.Type() {
		case model.TypeMap:
			if err := data.RangeMap(func(key string, v *model.Value) error {
				match, err := processValue(ctx, v, options)
				if err != nil {
					return err
				}

				if match {
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
				match, err := processValue(ctx, v, options)
				if err != nil {
					return err
				}

				if match {
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
