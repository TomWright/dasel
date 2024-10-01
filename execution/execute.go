package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
	"github.com/tomwright/dasel/v3/selector/lexer"
	"github.com/tomwright/dasel/v3/selector/parser"
)

func ExecuteSelector(selector string, value *model.Value) (*model.Value, error) {
	tokens, err := lexer.NewTokenizer(selector).Tokenize()
	if err != nil {
		return nil, fmt.Errorf("error tokenizing selector: %w", err)
	}

	expr, err := parser.NewParser(tokens).Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing selector: %w", err)
	}

	res, err := ExecuteAST(expr, value)
	if err != nil {
		return nil, fmt.Errorf("error executing selector: %w", err)
	}

	return res, nil
}

type expressionExecutor func(data *model.Value) (*model.Value, error)

func ExecuteAST(expr ast.Expr, value *model.Value) (*model.Value, error) {
	executor, err := exprExecutor(expr)
	if err != nil {
		return nil, fmt.Errorf("error evaluating expression: %w", err)
	}
	res, err := executor(value)
	if err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	return res, nil
}

func exprExecutor(expr ast.Expr) (expressionExecutor, error) {
	switch e := expr.(type) {
	case ast.BinaryExpr:
		return binaryExprExecutor(e)
	case ast.CallExpr:
		return callExprExecutor(e)
	case ast.ChainedExpr:
		return chainedExprExecutor(e)
	case ast.SpreadExpr:
		return spreadExprExecutor()
	case ast.RangeExpr:
		return rangeExprExecutor(e)
	case ast.IndexExpr:
		return indexExprExecutor(e)
	case ast.PropertyExpr:
		return propertyExprExecutor(e)
	case ast.VariableExpr:
		return variableExprExecutor(e)
	case ast.NumberIntExpr:
		return numberIntExprExecutor(e)
	case ast.NumberFloatExpr:
		return numberFloatExprExecutor(e)
	case ast.StringExpr:
		return stringExprExecutor(e)
	case ast.BoolExpr:
		return boolExprExecutor(e)
	case ast.ObjectExpr:
		return objectExprExecutor(e)
	case ast.MapExpr:
		return mapExprExecutor(e)
	default:
		return nil, fmt.Errorf("unhandled expression type: %T", e)
	}
}

func binaryExprExecutor(e ast.BinaryExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		left, err := ExecuteAST(e.Left, data)
		if err != nil {
			return nil, fmt.Errorf("error evaluating left expression: %w", err)
		}
		right, err := ExecuteAST(e.Right, data)
		if err != nil {
			return nil, fmt.Errorf("error evaluating right expression: %w", err)
		}

		switch e.Operator.Kind {
		case lexer.Plus:
			return left.Add(right)
		case lexer.Dash:
			return left.Subtract(right)
		case lexer.Star:
			return left.Multiply(right)
		case lexer.Slash:
			return left.Divide(right)
		case lexer.Percent:
			return left.Modulo(right)
		case lexer.GreaterThan:
			return left.GreaterThan(right)
		case lexer.GreaterThanOrEqual:
			return left.GreaterThanOrEqual(right)
		case lexer.LessThan:
			return left.LessThan(right)
		case lexer.LessThanOrEqual:
			return left.LessThanOrEqual(right)
		case lexer.Equal:
			return left.Equal(right)
		case lexer.NotEqual:
			return left.NotEqual(right)
		default:
			return nil, fmt.Errorf("unhandled operator: %s", e.Operator.Value)
		}
	}, nil
}

func chainedExprExecutor(e ast.ChainedExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		for _, expr := range e.Exprs {
			res, err := ExecuteAST(expr, data)
			if err != nil {
				return nil, fmt.Errorf("error executing expression: %w", err)
			}
			data = res
		}
		return data, nil
	}, nil
}

func spreadExprExecutor() (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		s := model.NewSliceValue()

		switch {
		case data.IsSlice():
			v, err := data.SliceValue()
			if err != nil {
				return nil, fmt.Errorf("error getting slice value: %w", err)
			}
			for _, sv := range v {
				s.Append(model.NewValue(sv))
			}
		case data.IsMap():
			v, err := data.MapValue()
			if err != nil {
				return nil, fmt.Errorf("error getting map value: %w", err)
			}
			for _, kv := range v.KeyValues() {
				s.Append(model.NewValue(kv.Value))
			}
		default:
			return nil, fmt.Errorf("cannot spread on type %s", data.Type())
		}

		return s, nil
	}, nil
}

func rangeExprExecutor(e ast.RangeExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		panic("not implemented")
	}, nil
}

func indexExprExecutor(e ast.IndexExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		panic("not implemented")
	}, nil
}

func propertyExprExecutor(e ast.PropertyExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		if !data.IsMap() {
			return nil, fmt.Errorf("expected map, got %s", data.Type())
		}
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

func variableExprExecutor(e ast.VariableExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		varName := e.Name
		if varName == "this" {
			return data, nil
		}
		return nil, fmt.Errorf("variable %s not found", varName)
	}, nil
}
