package execution

import (
	"fmt"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/selector/ast"
)

func executeAssign(left *model.Value, right *model.Value, _ ast.BinaryExpr) (*model.Value, error) {
	err := left.Set(right)
	if err != nil {
		return nil, fmt.Errorf("error setting value: %w", err)
	}
	return right, nil
}
