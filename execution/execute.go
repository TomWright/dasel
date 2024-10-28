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

type expressionExecutor func(data *model.Value) (*model.Value, error)

// ExecuteAST executes the given AST with the given input.
func ExecuteAST(expr ast.Expr, value *model.Value, options *Options) (*model.Value, error) {
	if expr == nil {
		return value, nil
	}

	executor, err := exprExecutor(options, expr)
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

	if err := value.RangeSlice(func(i int, v *model.Value) error {
		r, err := executor(v)
		if err != nil {
			return err
		}
		return res.Append(r)
	}); err != nil {
		return nil, fmt.Errorf("branch execution error: %w", err)
	}

	return res, nil
}

var unstableAstTypes = []reflect.Type{
	reflect.TypeFor[ast.BranchExpr](),
}

func exprExecutor(opts *Options, expr ast.Expr) (expressionExecutor, error) {
	if !opts.Unstable && (slices.Contains(unstableAstTypes, reflect.TypeOf(expr)) ||
		slices.Contains(unstableAstTypes, reflect.ValueOf(expr).Type())) {
		return nil, errors.New("unstable ast types are not enabled. to enable them use --unstable")
	}

	switch e := expr.(type) {
	case ast.BinaryExpr:
		return binaryExprExecutor(opts, e)
	case ast.UnaryExpr:
		return unaryExprExecutor(opts, e)
	case ast.CallExpr:
		return callExprExecutor(opts, e)
	case ast.ChainedExpr:
		return chainedExprExecutor(opts, e)
	case ast.SpreadExpr:
		return spreadExprExecutor()
	case ast.RangeExpr:
		return rangeExprExecutor(opts, e)
	case ast.IndexExpr:
		return indexExprExecutor(opts, e)
	case ast.PropertyExpr:
		return propertyExprExecutor(opts, e)
	case ast.VariableExpr:
		return variableExprExecutor(opts, e)
	case ast.NumberIntExpr:
		return numberIntExprExecutor(e)
	case ast.NumberFloatExpr:
		return numberFloatExprExecutor(e)
	case ast.StringExpr:
		return stringExprExecutor(e)
	case ast.BoolExpr:
		return boolExprExecutor(e)
	case ast.ObjectExpr:
		return objectExprExecutor(opts, e)
	case ast.MapExpr:
		return mapExprExecutor(opts, e)
	case ast.FilterExpr:
		return filterExprExecutor(opts, e)
	case ast.ConditionalExpr:
		return conditionalExprExecutor(opts, e)
	case ast.BranchExpr:
		return branchExprExecutor(opts, e)
	case ast.ArrayExpr:
		return arrayExprExecutor(opts, e)
	case ast.RegexExpr:
		// Noop
		return func(data *model.Value) (*model.Value, error) {
			return data, nil
		}, nil
	case ast.SortByExpr:
		return sortByExprExecutor(opts, e)
	case ast.NullExpr:
		return func(data *model.Value) (*model.Value, error) {
			return model.NewNullValue(), nil
		}, nil
	default:
		return nil, fmt.Errorf("unhandled expression type: %T", e)
	}
}

func chainedExprExecutor(options *Options, e ast.ChainedExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		for _, expr := range e.Exprs {

			if !data.IsBranch() {
				res, err := ExecuteAST(expr, data, options)
				if err != nil {
					return nil, fmt.Errorf("error executing expression: %w", err)
				}
				data = res
				continue
			}

			res := model.NewSliceValue()
			res.MarkAsBranch()
			if err := data.RangeSlice(func(i int, value *model.Value) error {
				r, err := ExecuteAST(expr, value, options)
				if err != nil {
					return fmt.Errorf("error executing expression: %w", err)
				}

				if err := res.Append(r); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return nil, err
			}
			data = res
		}
		return data, nil
	}, nil
}

func variableExprExecutor(opts *Options, e ast.VariableExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		varName := e.Name
		if varName == "this" {
			return data, nil
		}
		res, ok := opts.Vars[varName]
		if !ok {
			return nil, fmt.Errorf("variable %s not found", varName)
		}
		return res, nil
	}, nil
}
