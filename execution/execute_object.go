package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func objectExprExecutor(e ast.ObjectExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		obj := model.NewMapValue()
		for _, p := range e.Pairs {
			if ast.IsSpreadExpr(p.Key) && ast.IsSpreadExpr(p.Value) {
				if err := data.RangeMap(func(key string, value *model.Value) error {
					if err := obj.SetMapKey(key, value); err != nil {
						return fmt.Errorf("error setting map key: %w", err)
					}
					return nil
				}); err != nil {
					return nil, fmt.Errorf("error ranging map: %w", err)
				}
				continue
			}

			if ast.IsSpreadExpr(p.Key) {
				return nil, fmt.Errorf("cannot spread object key name")
			}

			key, err := ExecuteAST(p.Key, data)
			if err != nil {
				return nil, fmt.Errorf("error evaluating key: %w", err)
			}
			if !key.IsString() {
				return nil, fmt.Errorf("expected key to resolve to string, got %s", key.Type())
			}
			val, err := ExecuteAST(p.Value, data)
			if err != nil {
				return nil, fmt.Errorf("error evaluating value: %w", err)
			}
			keyStr, err := key.StringValue()
			if err := obj.SetMapKey(keyStr, val); err != nil {
				return nil, fmt.Errorf("error setting map key: %w", err)
			}
		}
		return obj, nil
	}, nil
}

func propertyExprExecutor(e ast.PropertyExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		key, err := ExecuteAST(e.Property, data)
		if err != nil {
			return nil, fmt.Errorf("error evaluating property: %w", err)
		}
		if !key.IsString() {
			return nil, fmt.Errorf("expected property to resolve to string, got %s", key.Type())
		}
		keyStr, err := key.StringValue()
		if err != nil {
			return nil, fmt.Errorf("error getting string value: %w", err)
		}
		return data.GetMapKey(keyStr)
	}, nil
}
