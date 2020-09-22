package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal"
)

func versionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the dasel version.",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(internal.Version)
		},
	}

	return cmd
}
