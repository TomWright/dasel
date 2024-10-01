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

func callSingleExecutor(f singleResponseFunc, argsE ast.Expressions) (expressionExecutor, error) {
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

func callMultiExecutor(f multiResponseFunc, argsE ast.Expressions) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		panic("multi response functions are not supported")
		//args, err := prepareArgs(data, argsE)
		//if err != nil {
		//	return nil, fmt.Errorf("error preparing arguments: %w", err)
		//}

		//res, err := f(data, args)
		//if err != nil {
		//	return nil, fmt.Errorf("error executing function: %w", err)
		//}

		//return res, nil
	}, nil
}

func callExprExecutor(e ast.CallExpr) (expressionExecutor, error) {
	if f, ok := singleResponseFuncLookup[e.Function]; ok {
		return callSingleExecutor(f, e.Args)
	}
	if f, ok := multiResponseFuncLookup[e.Function]; ok {
		return callMultiExecutor(f, e.Args)
	}

	return nil, fmt.Errorf("unknown function: %q", e.Function)
}
