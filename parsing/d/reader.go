package json

import (
	"context"
	"fmt"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// Dasel represents the dasel format.
	Dasel parsing.Format = "dasel"
)

var _ parsing.Reader = (*daselReader)(nil)

func init() {
	parsing.RegisterReader(Dasel, newDaselReader)
}

func newDaselReader(options parsing.ReaderOptions) (parsing.Reader, error) {
	return &daselReader{}, nil
}

type daselReader struct {
}

func (dr *daselReader) Read(in []byte) (*model.Value, error) {
	if len(in) == 0 {
		return model.NewNullValue(), nil
	}
	out, err := execution.ExecuteSelector(context.Background(), string(in), model.NewNullValue(), execution.NewOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to read value: %w", err)
	}
	return out, nil
}
