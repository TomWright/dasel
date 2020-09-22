package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal"
)

// RootCMD is the root command for use with cobra.
var RootCMD = &cobra.Command{
	Use:   "dasel",
	Short: "Query and modify data structures using selector strings.",
}

func init() {
	RootCMD.Version = internal.Version
	RootCMD.AddCommand(
		selectCommand(),
	)
}
