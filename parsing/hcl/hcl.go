package hcl

import (
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// HCL represents the hcl2 file format.
	HCL parsing.Format = "hcl"
)

var _ parsing.Reader = (*hclReader)(nil)
var _ parsing.Writer = (*hclWriter)(nil)

func init() {
	parsing.RegisterReader(HCL, newHCLReader)
	// HCL writer is not implemented yet
	//parsing.RegisterWriter(HCL, newHCLWriter)
}
