package command

import (
	"github.com/spf13/cobra"
)

func putBoolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bool -f <file> -s <selector> <value>",
		Short: "Update a bool property in the given document.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenericPutCommand(genericPutOptions{
				ValueType: "bool",
				Init:      getGenericInit(cmd, args),
			}, cmd)
		},
	}

	return cmd
}
