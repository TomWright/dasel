package execution

import (
	"context"
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"io"
	"os"
)

// FuncReadFile reads the given filepath at runtime.
var FuncReadFile = NewFunc(
	"readFile",
	func(ctx context.Context, data *model.Value, args model.Values) (*model.Value, error) {
		filepath, err := args[0].StringValue()
		if err != nil {
			return nil, fmt.Errorf("readFile: %w", err)
		}

		f, err := os.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("readFile: %w", err)
		}
		defer func() {
			_ = f.Close()
		}()

		fileBytes, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("readFile: %w", err)
		}

		return model.NewStringValue(string(fileBytes)), nil
	},
	ValidateArgsExactly(1),
)
