package execution

import (
	"github.com/tomwright/dasel/v3/model"
)

// FuncMax is a function that returns the highest number.
var FuncMax = NewFunc(
	"max",
	func(data *model.Value, args model.Values) (*model.Value, error) {
		res := model.NewNullValue()
		for _, arg := range args {
			if res.IsNull() {
				res = arg
				continue
			}
			gt, err := arg.GreaterThan(res)
			if err != nil {
				return nil, err
			}
			gtBool, err := gt.BoolValue()
			if err != nil {
				return nil, err
			}
			if gtBool {
				res = arg
			}
		}
		return res, nil
	},
	ValidateArgsMin(1),
)
