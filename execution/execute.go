package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func ExecuteSelector(selectorStr string, value *model.Value) (*model.Value, error) {
	if selectorStr == "" {
		return value, nil
	}

	expr, err := selector.Parse(selectorStr)
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
	if expr == nil {
		return value, nil
	}

	executor, err := exprExecutor(expr)
	if err != nil {
		return nil, fmt.Errorf("error evaluating expression: %w", err)
	}

	if !value.IsBranch() {
		res, err := executor(value)
		if err != nil {
			return nil, fmt.Errorf("execution error: %w", err)
		}
		return res, nil
	}

	res := model.NewSliceValue()
	res.MarkAsBranch()

	if err := value.RangeSlice(func(i int, value *model.Value) error {
		r, err := executor(value)
		if err != nil {
			return err
		}
		return res.Append(r)
	}); err != nil {
		return nil, fmt.Errorf("branch execution error: %w", err)
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
	case ast.ConditionalExpr:
		return conditionalExprExecutor(e)
	case ast.BranchExpr:
		return branchExprExecutor(e)
	case ast.ArrayExpr:
		return arrayExprExecutor(e)
	default:
		return nil, fmt.Errorf("unhandled expression type: %T", e)
	}
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

func variableExprExecutor(e ast.VariableExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		varName := e.Name
		if varName == "this" {
			return data, nil
		}
		return nil, fmt.Errorf("variable %s not found", varName)
	}, nil
}
