package execution

import (
	"context"
	"errors"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func recursiveDescentExprExecutor2(e ast.RecursiveDescentExpr) (expressionExecutor, error) {
	var doSearch func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error)
	findValue := func(ctx context.Context, options *Options, v *model.Value) (*model.Value, error) {
		property, err := ExecuteAST(ctx, e.Expr, v, options)
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
			return nil, err
		}
		return property, nil
	}
	doSearch = func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error) {
		res := make([]*model.Value, 0)

		switch data.Type() {
		case model.TypeMap:
			if err := data.RangeMap(func(key string, v *model.Value) error {
				if v.IsScalar() {
					if e.IsWildcard {
						res = append(res, v)
					}
				} else {
					if !e.IsWildcard {
						property, err := findValue(ctx, options, v)
						if err != nil {
							return err
						}
						if property != nil {
							res = append(res, property)
						}
					}

					gotNext, err := doSearch(ctx, options, v)
					if err != nil {
						return err
					}
					res = append(res, gotNext...)
				}
				return nil
			}); err != nil {
				return nil, err
			}
		case model.TypeSlice:
			if err := data.RangeSlice(func(i int, v *model.Value) error {
				if v.IsScalar() {
					if e.IsWildcard {
						res = append(res, v)
					}
				} else {
					if !e.IsWildcard {
						property, err := findValue(ctx, options, v)
						if err != nil {
							return err
						}
						if property != nil {
							res = append(res, property)
						}
					}

					gotNext, err := doSearch(ctx, options, v)
					if err != nil {
						return err
					}
					res = append(res, gotNext...)
				}
				return nil
			}); err != nil {
				return nil, err
			}
		}

		return res, nil
	}

	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "recursiveDescentExpr")
		matches := model.NewSliceValue()

		found, err := doSearch(ctx, options, data)
		if err != nil {
			return nil, err
		}

		for _, f := range found {
			// We purposely wrap the value here to ensure any downstream changes are applied to the root.
			if err := matches.Append(model.NewValue(f)); err != nil {
				return nil, err
			}
		}

		return matches, nil
	}, nil
}
