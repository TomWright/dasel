package command

import (
	"github.com/spf13/cobra"
)

func putStringCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "string -f <file> -s <selector> <value>",
		Short: "Update a string property in the given document.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenericPutCommand(genericPutOptions{
				ValueType: "string",
				Init:      getGenericInit(cmd, args),
			}, cmd)
		},
	}
}
