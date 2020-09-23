package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func putBoolCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bool -f <file> -s <selector> <value>",
		Short: "Update a bool property in the given file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flag("file").Value.String()
			outFlag := cmd.Flag("out").Value.String()
			parserFlag := cmd.Flag("parser").Value.String()
			selectorFlag := cmd.Flag("selector").Value.String()

			parser, err := getParser(fileFlag, parserFlag)
			if err != nil {
				return err
			}
			rootNode, err := getRootNode(fileFlag, parser)
			if err != nil {
				return err
			}

			updateValue, err := parseValue(args[0], "bool")
			if err != nil {
				return err
			}

			if err := rootNode.Put(selectorFlag, updateValue); err != nil {
				return fmt.Errorf("could not put value: %w", err)
			}

			if err := writeNodeToOutput(rootNode, parser, fileFlag, outFlag); err != nil {
				return fmt.Errorf("could not write output: %w", err)
			}

			return nil
		},
	}

	return cmd
}
