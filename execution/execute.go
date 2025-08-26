package execution

import (
	"context"
	"errors"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector"
	"github.com/tomwright/dasel/v3/selector/ast"
	"os"
	"reflect"
	"slices"
)

// ExecuteSelector parses the selector and executes the resulting AST with the given input.
func ExecuteSelector(ctx context.Context, selectorStr string, value *model.Value, opts *Options) (*model.Value, error) {
	if selectorStr == "" {
		return value, nil
	}

	expr, err := selector.Parse(selectorStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing selector: %w", err)
	}

	res, err := ExecuteAST(ctx, expr, value, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing selector: %w", err)
	}

	return res, nil
}

type expressionExecutor func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error)

// ExecuteAST executes the given AST with the given input.
func ExecuteAST(ctx context.Context, expr ast.Expr, value *model.Value, options *Options) (*model.Value, error) {
	if expr == nil {
		return value, nil
	}

	executorFn, err := exprExecutor(options, expr)
	if err != nil {
		return nil, fmt.Errorf("error evaluating expression %T: %w", expr, err)
	}

	executor := func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		options.Vars["this"] = data
		out, err := executorFn(ctx, options, data)
		if err != nil {
			return out, err
		}
		return out, nil
	}

	if !value.IsBranch() {
		res, err := executor(ctx, options, value)
		if err != nil {
			return nil, fmt.Errorf("execution error when processing %T: %w", expr, err)
		}
		return res, nil
	}

	res := model.NewSliceValue()
	res.MarkAsBranch()

	if err := value.RangeSlice(func(i int, v *model.Value) error {
		r, err := executor(ctx, options, v)
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
		return recursiveDescentExprExecutor2(e)
	case ast.ConditionalExpr:
		return conditionalExprExecutor(e)
	case ast.BranchExpr:
		return branchExprExecutor(e)
	case ast.ArrayExpr:
		return arrayExprExecutor(e)
	case ast.RegexExpr:
		// Noop
		return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
			ctx = WithExecutorID(ctx, "regexExpr")
			return data, nil
		}, nil
	case ast.SortByExpr:
		return sortByExprExecutor(e)
	case ast.NullExpr:
		return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
			ctx = WithExecutorID(ctx, "nullExpr")
			return model.NewNullValue(), nil
		}, nil
	default:
		return nil, fmt.Errorf("unhandled expression type: %T", e)
	}
}

func chainedExprExecutor(e ast.ChainedExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "chainedExpr")
		var curData = data
		for _, expr := range e.Exprs {
			res, err := ExecuteAST(ctx, expr, curData, options)
			if err != nil {
				return nil, fmt.Errorf("error executing expression: %w", err)
			}
			curData = res
		}
		return curData, nil
	}, nil
}

func variableExprExecutor(e ast.VariableExpr) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "variableExpr")
		varName := e.Name
		res, ok := options.Vars[varName]
		if ok {
			return res, nil
		}

		envVarValue := os.Getenv(varName)
		if envVarValue != "" {
			return model.NewStringValue(envVarValue), nil
		}

		return nil, fmt.Errorf("variable %s not found", varName)
	}, nil
}
