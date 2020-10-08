package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal"
)

// NewRootCMD returns the root command for use with cobra.
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dasel",
		Short: "Query and modify data structures using selector strings.",
	}
	cmd.Version = internal.Version
	cmd.AddCommand(
		selectCommand(),
		putCommand(),
	)
	return cmd
}
