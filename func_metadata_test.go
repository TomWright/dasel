package dasel

import (
	"testing"
)

func TestMetadataFunc(t *testing.T) {
	t.Run("Args", selectTestErr(
		"metadata()",
		map[string]interface{}{},
		&ErrUnexpectedFunctionArgs{
			Function: "metadata",
			Args:     []string{},
		}),
	)
}
