package execution

import (
	"github.com/tomwright/dasel/v3/model"
)

// FuncMin is a function that returns the smalled number.
var FuncMin = NewFunc(
	"min",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		res := model.NewNullValue()
		for _, arg := range args {
			if res.IsNull() {
				res = arg
				continue
			}
			lt, err := arg.LessThan(res)
			if err != nil {
				return nil, err
			}
			ltBool, err := lt.BoolValue()
			if err != nil {
				return nil, err
			}
			if ltBool {
				res = arg
			}
		}
		return res, nil
	},
	ValidateArgsMin(1),
)
