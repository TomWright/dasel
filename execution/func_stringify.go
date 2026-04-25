package execution

import (
	"bytes"
	"context"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// FuncStringify serializes a structured value into a format string.
var FuncStringify = NewFunc(
	"stringify",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		formatStr, err := args[0].StringValue()
		if err != nil {
			return nil, err
		}
		format := parsing.Format(formatStr)

		input := data
		if len(args) == 2 {
			input = args[1]
		}

		writer, err := format.NewWriter(parsing.DefaultWriterOptions())
		if err != nil {
			return nil, err
		}

		b, err := writer.Write(input)
		if err != nil {
			return nil, err
		}

		b = bytes.TrimSuffix(b, []byte("\n"))

		return model.NewStringValue(string(b)), nil
	},
	ValidateArgsMinMax(1, 2),
)
