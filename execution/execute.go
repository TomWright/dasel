package execution

import (
	"errors"
	"fmt"
	"reflect"
	"slices"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector"
	"github.com/tomwright/dasel/v3/selector/ast"
)

// ExecuteSelector parses the selector and executes the resulting AST with the given input.
func ExecuteSelector(selectorStr string, value *model.Value, opts *Options) (*model.Value, error) {
	if selectorStr == "" {
		return value, nil
	}

	expr, err := selector.Parse(selectorStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing selector: %w", err)
	}

	res, err := ExecuteAST(expr, value, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing selector: %w", err)
	}

	return res, nil
}

type expressionExecutor func(options *Options, data *model.Value) (*model.Value, error)

// ExecuteAST executes the given AST with the given input.
func ExecuteAST(expr ast.Expr, value *model.Value, options *Options) (*model.Value, error) {
	if expr == nil {
		return value, nil
	}

	executor, err := exprExecutor(options, expr)
	if err != nil {
		return nil, fmt.Errorf("error evaluating expression %T: %w", expr, err)
	}

	if !value.IsBranch() {
		options.Vars["this"] = value
		res, err := executor(options, value)
		if err != nil {
			return nil, fmt.Errorf("execution error when processing %T: %w", expr, err)
		}
		return res, nil
	}

	res := model.NewSliceValue()
	res.MarkAsBranch()

	if err := value.RangeSlice(func(i int, v *model.Value) error {
		options.Vars["this"] = v
		r, err := executor(options, v)
		if err != nil {
			return err
		}
		if r.IsIgnore() {
			return nil
		}
		return res.Append(r)
	}); err != nil {
		return nil, fmt.Errorf("branch execution error when processing %T: %w", expr, err)
	}

	return res, nil
}

var unstableAstTypes = []reflect.Type{
	reflect.TypeFor[ast.BranchExpr](),
}

func exprExecutor(options *Options, expr ast.Expr) (expressionExecutor, error) {
	if !options.Unstable && (slices.Contains(unstableAstTypes, reflect.TypeOf(expr)) ||
		slices.Contains(unstableAstTypes, reflect.ValueOf(expr).Type())) {
		return nil, errors.New("unstable ast types are not enabled. to enable them use --unstable")
	}

	switch e := expr.(type) {
	case ast.BinaryExpr:
		return binaryExprExecutor(e)
	case ast.UnaryExpr:
		return unaryExprExecutor(e)
	case ast.CallExpr:
		return callExprExecutor(options, e)
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
	case ast.EachExpr:
		return eachExprExecutor(e)
	case ast.FilterExpr:
		return filterExprExecutor(e)
	case ast.SearchExpr:
		return searchExprExecutor(e)
	case ast.RecursiveDescentExpr:
		return recursiveDescentExprExecutor(e)
	case ast.ConditionalExpr:
		return conditionalExprExecutor(e)
	case ast.BranchExpr:
		return branchExprExecutor(e)
	case ast.ArrayExpr:
		return arrayExprExecutor(e)
	case ast.RegexExpr:
		// Noop
		return func(options *Options, data *model.Value) (*model.Value, error) {
			return data, nil
		}, nil
	case ast.SortByExpr:
		return sortByExprExecutor(e)
	case ast.NullExpr:
		return func(options *Options, data *model.Value) (*model.Value, error) {
			return model.NewNullValue(), nil
		}, nil
	default:
		return nil, fmt.Errorf("unhandled expression type: %T", e)
	}
}

func chainedExprExecutor(e ast.ChainedExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		for _, expr := range e.Exprs {
			res, err := ExecuteAST(expr, data, options)
			if err != nil {
				return nil, fmt.Errorf("error executing expression: %w", err)
			}
			data = res
		}
		return data, nil
	}, nil
}

func variableExprExecutor(e ast.VariableExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		varName := e.Name
		res, ok := options.Vars[varName]
		if !ok {
			return nil, fmt.Errorf("variable %s not found", varName)
		}
		return res, nil
	}, nil
}
