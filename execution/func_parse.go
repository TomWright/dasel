package execution

import (
	"context"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// FuncParse parses the given data at runtime.
var FuncParse = NewFunc(
	"parse",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		var format parsing.Format
		var content []byte
		{
			strVal, err := args[0].StringValue()
			if err != nil {
				return nil, err
			}
			format = parsing.Format(strVal)
		}
		{
			strVal, err := args[1].StringValue()
			if err != nil {
				return nil, err
			}
			content = []byte(strVal)
		}

		reader, err := format.NewReader(parsing.DefaultReaderOptions())
		if err != nil {
			return nil, err
		}

		doc, err := reader.Read(content)
		if err != nil {
			return nil, err
		}

		return doc, nil
	},
	ValidateArgsExactly(2),
)
