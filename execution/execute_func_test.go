package execution_test

import (
	"testing"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
)

func TestFunc(t *testing.T) {
	returnInputData := execution.NewFunc(
		"returnInputData",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			return data, nil
		},
		execution.ValidateArgsExactly(0),
	)

	returnFirstArg := execution.NewFunc(
		"returnFirstArg",
		func(data *model.Value, args model.Values) (*model.Value, error) {
			return args[0], nil
		},
		execution.ValidateArgsExactly(1),
	)

	funcs := execution.NewFuncCollection(
		returnInputData,
		returnFirstArg,
	)

	opts := []execution.ExecuteOptionFn{
		func(options *execution.Options) {
			options.Funcs = funcs
		},
	}

	t.Run("returnInputData", testCase{
		s:    `1.returnInputData()`,
		out:  model.NewIntValue(1),
		opts: opts,
	}.run)

	t.Run("returnFirstArg", testCase{
		s:    `1.returnFirstArg(2)`,
		out:  model.NewIntValue(2),
		opts: opts,
	}.run)
}
