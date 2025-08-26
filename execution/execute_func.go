package execution

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func prepareArgs(ctx context.Context, opts *Options, data *model.Value, argsE ast.Expressions) (model.Values, error) {
	args := make(model.Values, 0)
	for i, arg := range argsE {
		res, err := ExecuteAST(ctx, arg, data, opts)
		if err != nil {
			return nil, fmt.Errorf("error evaluating argument %d: %w", i, err)
		}

		argVals, err := prepareSpreadValues(res)
		if err != nil {
			return nil, fmt.Errorf("error handling spread values: %w", err)
		}

		args = append(args, argVals...)
	}
	return args, nil
}

func callFnExecutor(f FuncFn, argsE ast.Expressions) (expressionExecutor, error) {
	return func(ctx context.Context, options *Options, data *model.Value) (*model.Value, error) {
		ctx = WithExecutorID(ctx, "callFnExpr")
		args, err := prepareArgs(ctx, options, data, argsE)
		if err != nil {
			return nil, fmt.Errorf("error preparing arguments: %w", err)
		}

		res, err := f(ctx, data, args)
		if err != nil {
			return nil, fmt.Errorf("error executing function: %w", err)
		}

		return res, nil
	}, nil
}

var unstableFuncs = []string{
	"ignore",
}

func callExprExecutor(options *Options, e ast.CallExpr) (expressionExecutor, error) {
	if !options.Unstable && (slices.Contains(unstableFuncs, e.Function)) {
		return nil, errors.New("unstable function are not enabled. to enable them use --unstable")
	}
	if f, ok := options.Funcs.Get(e.Function); ok {
		res, err := callFnExecutor(f, e.Args)
		if err != nil {
			return nil, fmt.Errorf("error executing function %q: %w", e.Function, err)
		}
		return res, nil
	}

	return nil, fmt.Errorf("unknown function: %q", e.Function)
}
