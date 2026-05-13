package kdl

import "github.com/tomwright/dasel/v3/parsing"

// KDL represents the KDL file format.
const KDL parsing.Format = "kdl"

func init() {
	parsing.RegisterReader(KDL, newKDLReader)
	parsing.RegisterWriter(KDL, newKDLWriter)
}
