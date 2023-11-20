package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func manCommand(root *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "man -o <dir>",
		Short: "Generate manual pages for all dasel subcommands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return manRunE(cmd, root)
		},
	}

	cmd.Flags().StringP("output-directory", "o", ".", "The directory in which man pages will be created")

	return cmd
}

func manRunE(cmd *cobra.Command, root *cobra.Command) error {
	outputDirectory, _ := cmd.Flags().GetString("output-directory")

	return doc.GenManTree(root, nil, outputDirectory)
}
