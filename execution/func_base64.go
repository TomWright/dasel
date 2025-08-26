package execution

import (
	"context"
	"encoding/base64"

	"github.com/tomwright/dasel/v3/model"
)

// FuncBase64Encode base64 encodes the given value.
var FuncBase64Encode = NewFunc(
	"base64e",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		arg := args[0]
		strVal, err := arg.StringValue()
		if err != nil {
			return nil, err
		}
		out := base64.StdEncoding.EncodeToString([]byte(strVal))
		return model.NewStringValue(out), nil
	},
	ValidateArgsExactly(1),
)

// FuncBase64Decode base64 decodes the given value.
var FuncBase64Decode = NewFunc(
	"base64d",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		arg := args[0]
		strVal, err := arg.StringValue()
		if err != nil {
			return nil, err
		}
		out, err := base64.StdEncoding.DecodeString(strVal)
		if err != nil {
			return nil, err
		}
		return model.NewStringValue(string(out)), nil
	},
	ValidateArgsExactly(1),
)
