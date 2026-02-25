package execution

import (
	"context"
	"strings"

	"github.com/tomwright/dasel/v3/model"
)

// FuncReplace is a function that replaces all occurrences of a substring with another string.
var FuncReplace = NewFunc(
	"replace",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		inputData := data
		if len(args)%2 != 0 {
			inputData = args[0]
			args = args[1:]
		}

		argStrings := make([]string, len(args))
		for i, arg := range args {
			s, err := arg.StringValue()
			if err != nil {
				return nil, err
			}
			argStrings[i] = s
		}
		replacer := strings.NewReplacer(argStrings...)

		inputString, err := inputData.StringValue()
		if err != nil {
			return nil, err
		}

		outputString := replacer.Replace(inputString)

		return model.NewStringValue(outputString), nil
	},
	ValidateArgsMin(2),
)
