package execution

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func doesValueMatchRecursiveDescentKey(opts *Options, data *model.Value, e ast.RecursiveDescentExpr) (*model.Value, error) {
	if e.IsWildcard {
		if data.IsScalar() {
			return data, nil
		}
		return nil, nil
	}

	var key *model.Value

	var expr ast.Expr

	switch exprT := e.Expr.(type) {
	case ast.PropertyExpr:
		expr = exprT.Property
	case ast.IndexExpr:
		expr = exprT.Index
	default:
		expr = e.Expr
	}

	key, err := ExecuteAST(expr, data, opts)
	if err != nil {
		return nil, err
	}

	switch key.Type() {
	case model.TypeString:
		keyStr, err := key.StringValue()
		if err != nil {
			return nil, err
		}
		if !data.IsMap() {
			return nil, nil
		}
		exists, err := data.MapKeyExists(keyStr)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, nil
		}
		return data.GetMapKey(keyStr)

	case model.TypeInt:
		keyInt, err := key.IntValue()
		if err != nil {
			return nil, err
		}
		if !data.IsSlice() {
			return nil, nil
		}
		sliceSize, err := data.SliceLen()
		if err != nil {
			return nil, err
		}
		if keyInt >= 0 && keyInt < int64(sliceSize) {
			res, err := data.GetSliceIndex(int(keyInt))
			return res, err
		}
		return nil, nil
	default:
		// TODO : Do we need to handle variable lookup?
		return nil, fmt.Errorf("unexpected recursive descent key type: %v", key.Type())
	}
}

func recursiveDescentExprExecutor(e ast.RecursiveDescentExpr) (expressionExecutor, error) {
	var recurseTree func(options *Options, data *model.Value) ([]*model.Value, error)

	recurseTree = func(options *Options, data *model.Value) ([]*model.Value, error) {
		res := make([]*model.Value, 0)

		switch data.Type() {
		case model.TypeMap:
			if err := data.RangeMap(func(key string, v *model.Value) error {
				appendValue, err := doesValueMatchRecursiveDescentKey(options, v, e)
				if err != nil {
					return err
				}
				if appendValue != nil {
					res = append(res, appendValue)
				}

				gotNext, err := recurseTree(options, v)
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

				appendValue, err := doesValueMatchRecursiveDescentKey(options, v, e)
				if err != nil {
					return err
				}
				if appendValue != nil {
					res = append(res, appendValue)
				}

				gotNext, err := recurseTree(options, v)
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

	return func(options *Options, data *model.Value) (*model.Value, error) {
		matches := model.NewSliceValue()

		found, err := recurseTree(options, data)
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
