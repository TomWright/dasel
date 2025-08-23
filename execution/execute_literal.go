package execution

import (
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func numberIntExprExecutor(e ast.NumberIntExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		return model.NewIntValue(e.Value), nil
	}, nil
}

func numberFloatExprExecutor(e ast.NumberFloatExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		return model.NewFloatValue(e.Value), nil
	}, nil
}

func stringExprExecutor(e ast.StringExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		return model.NewStringValue(e.Value), nil
	}, nil
}

func boolExprExecutor(e ast.BoolExpr) (expressionExecutor, error) {
	return func(options *Options, data *model.Value) (*model.Value, error) {
		return model.NewBoolValue(e.Value), nil
	}, nil
}
