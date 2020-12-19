package command

import (
	"github.com/spf13/cobra"
)

func putIntCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "int -f <file> -s <selector> <value>",
		Short: "Update an int property in the given document.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenericPutCommand(genericPutOptions{
				ValueType: "int",
				Init:      getGenericInit(cmd, args),
			}, cmd)
		},
	}

	return cmd
}
