package kdl

import (
	"github.com/tomwright/dasel/v3/parsing"
	"github.com/tomwright/dasel/v3/parsing/kdl/internal"
)

// KDL represents the KDL file format.
const KDL parsing.Format = "kdl"

// ErrKDLMaxDepthExceeded is returned when KDL children nesting depth exceeds the
// parser limit. The parser is mutually recursive over children blocks, so the
// limit stops deeply nested input from overflowing the stack.
var ErrKDLMaxDepthExceeded = internal.ErrMaxDepthExceeded

func init() {
	parsing.RegisterReader(KDL, newKDLReader)
	parsing.RegisterWriter(KDL, newKDLWriter)
}
