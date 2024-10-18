package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func prepareArgs(data *model.Value, argsE ast.Expressions) (model.Values, error) {
	args := make(model.Values, 0)
	for i, arg := range argsE {
		res, err := ExecuteAST(arg, data)
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
	return func(data *model.Value) (*model.Value, error) {
		args, err := prepareArgs(data, argsE)
		if err != nil {
			return nil, fmt.Errorf("error preparing arguments: %w", err)
		}

		res, err := f(data, args)
		if err != nil {
			return nil, fmt.Errorf("error executing function: %w", err)
		}

		return res, nil
	}, nil
}

func callExprExecutor(opts *Options, e ast.CallExpr) (expressionExecutor, error) {
	if f, ok := opts.Funcs.Get(e.Function); ok {
		res, err := callFnExecutor(f, e.Args)
		if err != nil {
			return nil, fmt.Errorf("error executing function %q: %w", e.Function, err)
		}
		return res, nil
	}

	return nil, fmt.Errorf("unknown function: %q", e.Function)
}
