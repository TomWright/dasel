package ini

import (
	"github.com/tomwright/dasel/v3/parsing"
)

const (
	// INI represents the ini file format.
	INI parsing.Format = "ini"
)

func init() {
	parsing.RegisterReader(INI, newINIReader)
	parsing.RegisterWriter(INI, newINIWriter)
}
