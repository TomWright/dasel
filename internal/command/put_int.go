package command

import (
	"github.com/spf13/cobra"
)

func putIntCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "int -f <file> -s <selector> <value>",
		Short: "Update an int property in the given document.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenericPutCommand(genericPutOptions{
				ValueType: "int",
				Init:      getGenericInit(cmd, args),
			}, cmd)
		},
	}
}
