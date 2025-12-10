package toml

import (
	"github.com/tomwright/dasel/v3/parsing"
)

// TODO : Implement using https://github.com/pelletier/go-toml/blob/v2/unstable/ast.go

// TOML represents the TOML file format.
const TOML parsing.Format = "toml"

func init() {
	parsing.RegisterReader(TOML, newTOMLReader)
	parsing.RegisterWriter(TOML, newTOMLWriter)
}
