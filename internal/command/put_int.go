package command

import (
	"github.com/spf13/cobra"
)

func putIntCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "int -f <file> -s <selector> <value>",
		Short: "Update an int property in the given file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenericPutCommand(genericPutOptions{
				Value:     args[0],
				ValueType: "int",
				Init:      getGenericInit(cmd),
			}, cmd)
		},
	}

	return cmd
}
