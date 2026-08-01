package plain

import (
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// Plain represents plain, unquoted scalar output.
	Plain parsing.Format = "plain"
)

func init() {
	// Write only. Plain output carries no structure, so there is nothing to read back.
	parsing.RegisterWriter(Plain, newPlainWriter)
}
