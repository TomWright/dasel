package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v2/internal"
)

// NewRootCMD returns the root command for use with cobra.
func NewRootCMD() *cobra.Command {
	selectCmd := selectCommand()
	selectCmd.SilenceErrors = true
	selectCmd.SilenceUsage = true
	selectCmd.Version = internal.Version

	selectCmd.AddCommand(
		putCommand(),
		deleteCommand(),
		validateCommand(),
	)

	return selectCmd
}
