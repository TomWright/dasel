package execution

import (
	"fmt"

	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func arrayExprExecutor(opts *Options, e ast.ArrayExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		res := model.NewSliceValue()

		for _, expr := range e.Exprs {
			el, err := ExecuteAST(expr, data, opts)
			if err != nil {
				return nil, err
			}
			if err := res.Append(el); err != nil {
				return nil, err
			}
		}

		return res, nil
	}, nil
}

func rangeExprExecutor(opts *Options, e ast.RangeExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		var start, end int64 = -1, -1
		if e.Start != nil {
			startE, err := ExecuteAST(e.Start, data, opts)
			if err != nil {
				return nil, fmt.Errorf("error evaluating start expression: %w", err)
			}

			start, err = startE.IntValue()
			if err != nil {
				return nil, fmt.Errorf("error getting start int value: %w", err)
			}
		}

		if e.End != nil {
			endE, err := ExecuteAST(e.End, data, opts)
			if err != nil {
				return nil, fmt.Errorf("error evaluating end expression: %w", err)
			}

			end, err = endE.IntValue()
			if err != nil {
				return nil, fmt.Errorf("error getting end int value: %w", err)
			}
		}

		res, err := data.SliceIndexRange(int(start), int(end))
		if err != nil {
			return nil, fmt.Errorf("error getting slice index range: %w", err)
		}

		return res, nil
	}, nil
}

func indexExprExecutor(opts *Options, e ast.IndexExpr) (expressionExecutor, error) {
	return func(data *model.Value) (*model.Value, error) {
		indexE, err := ExecuteAST(e.Index, data, opts)
		if err != nil {
			return nil, fmt.Errorf("error evaluating index expression: %w", err)
		}

		index, err := indexE.IntValue()
		if err != nil {
			return nil, fmt.Errorf("error getting index int value: %w", err)
		}

		return data.GetSliceIndex(int(index))
	}, nil
}