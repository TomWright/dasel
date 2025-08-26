package execution

import (
	"context"
	"errors"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

//func doesValueMatchRecursiveDescentKey(ctx context.Context, opts *Options, data *model.Value, e ast.RecursiveDescentExpr) (*model.Value, error) {
//	if e.IsWildcard {
//		if data.IsScalar() {
//			return data, nil
//		}
//		return nil, nil
//	}
//
//	property, err := ExecuteAST(ctx, e.Expr, data, opts)
//	if err != nil {
//		handleErrs := []any{
//			model.ErrIncompatibleTypes{},
//			model.ErrUnexpectedType{},
//			model.ErrUnexpectedTypes{},
//			model.SliceIndexOutOfRange{},
//			model.MapKeyNotFound{},
//		}
//		for _, e := range handleErrs {
//			if errors.As(err, &e) {
//				err = nil
//				break
//			}
//		}
//
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return property, nil
//}
//
//func recursiveDescentExprExecutor(e ast.RecursiveDescentExpr) (expressionExecutor, error) {
//	var recurseTree func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error)
//
//	recurseTree = func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error) {
//		res := make([]*model.Value, 0)
//
//		switch data.Type() {
//		case model.TypeMap:
//			if err := data.RangeMap(func(key string, v *model.Value) error {
//				appendValue, err := doesValueMatchRecursiveDescentKey(ctx, options, v, e)
//				if err != nil {
//					return err
//				}
//				if appendValue != nil {
//					res = append(res, appendValue)
//				}
//
//				gotNext, err := recurseTree(ctx, options, v)
//				if err != nil {
//					return err
//				}
//				res = append(res, gotNext...)
//
//				return nil
//			}); err != nil {
//				return nil, err
//			}
//		case model.TypeSlice:
//			if err := data.RangeSlice(func(i int, v *model.Value) error {
//
//				appendValue, err := doesValueMatchRecursiveDescentKey(ctx, options, v, e)
//				if err != nil {
//					return err
//				}
//				if appendValue != nil {
//					res = append(res, appendValue)
//				}
//
//				gotNext, err := recurseTree(ctx, options, v)
//				if err != nil {
//					return err
//				}
//				res = append(res, gotNext...)
//
//				return nil
//			}); err != nil {
//				return nil, err
//			}
//		}
//
//		return res, nil
//	}
//
//	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
//		ctx = WithExecutorID(ctx, "recursiveDescentExpr")
//		matches := model.NewSliceValue()
//
//		found, err := recurseTree(ctx, options, data)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, f := range found {
//			if f.Value.CanAddr() {
//				f.Value = f.Value.Addr()
//			} else {
//				ptr := reflect.New(f.Value.Type())
//				ptr.Elem().Set(f.Value)
//				f.Value = ptr
//			}
//			if err := matches.Append(f); err != nil {
//				return nil, err
//			}
//		}
//
//		before := data.String()
//
//		beforeIDs := make([]string, 0)
//		afterIDs := make([]string, 0)
//
//		if err := matches.RangeSlice(func(i int, v *model.Value) error {
//			beforeIDs = append(beforeIDs, v.UUID())
//			//This doesn't work
//			if err := v.Set(model.NewIntValue(5)); err != nil {
//				return err
//			}
//			return nil
//		}); err != nil {
//			panic(err)
//		}
//
//		after := data.String()
//		for _, f := range found {
//			afterIDs = append(afterIDs, f.UUID())
//			// This works
//			if err := f.Set(model.NewIntValue(5)); err != nil {
//				panic(err)
//			}
//		}
//		after2 := data.String()
//
//		fmt.Println("before:", before, after, after2)
//
//		return matches, nil
//	}, nil
//}

func recursiveDescentExprExecutor2(e ast.RecursiveDescentExpr) (expressionExecutor, error) {
	var doSearch func(ctx context.Context, options *Options, data *model.Value) ([]*model.Value, error)
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
