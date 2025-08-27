package execution

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func objectExprExecutor(e ast.ObjectExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "objectExpr")
		obj := model.NewMapValue()
		for _, p := range e.Pairs {

			if ast.IsType[ast.SpreadExpr](p.Key) {
				var val *model.Value
				var err error
				if p.Value != nil {
					// We need to spread the resulting value.
					val, err = ExecuteAST(ctx, p.Value, data, options)
					if err != nil {
						return nil, fmt.Errorf("error evaluating spread values: %w", err)
					}
				} else {
					val = data
				}

				if err := val.RangeMap(func(key string, value *model.Value) error {
					if err := obj.SetMapKey(key, value); err != nil {
						return fmt.Errorf("error setting map key: %w", err)
					}
					return nil
				}); err != nil {
					return nil, fmt.Errorf("error spreading into object: %w", err)
				}
				continue
			}

			key, err := ExecuteAST(ctx, p.Key, data, options)
			if err != nil {
				return nil, fmt.Errorf("error evaluating key: %w", err)
			}
			if !key.IsString() {
				return nil, fmt.Errorf("expected key to resolve to string, got %s", key.Type())
			}

			val, err := ExecuteAST(ctx, p.Value, data, options)
			if err != nil {
				return nil, fmt.Errorf("error evaluating value: %w", err)
			}

			keyStr, err := key.StringValue()
			if err != nil {
				return nil, fmt.Errorf("error getting string value: %w", err)
			}
			if err := obj.SetMapKey(keyStr, val); err != nil {
				return nil, fmt.Errorf("error setting map key: %w", err)
			}
		}
		return obj, nil
	}, nil
}

func propertyExprExecutor(e ast.PropertyExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "propertyExpr")
		key, err := ExecuteAST(ctx, e.Property, data, options)
		if err != nil {
			return nil, fmt.Errorf("error evaluating property: %w", err)
		}
		switch {
		case key.IsString():
			keyStr, err := key.StringValue()
			if err != nil {
				return nil, fmt.Errorf("error getting string value: %w", err)
			}

			return data.GetMapKey(keyStr)
		case key.IsInt():
			keyInt, err := key.IntValue()
			if err != nil {
				return nil, fmt.Errorf("error getting int value: %w", err)
			}
			return data.GetSliceIndex(int(keyInt))
		default:
			return nil, fmt.Errorf("expected key to be a string or int, got %s", key.Type())
		}
	}, nil
}
