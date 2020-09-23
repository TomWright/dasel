package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
)

func runSelectCommand(fileFlag string, parserFlag string, selectorFlag string) error {
	parser, err := getParser(fileFlag, parserFlag)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(fileFlag, parser)
	if err != nil {
		return err
	}

	var res *dasel.Node
	if selectorFlag == "." {
		res = rootNode
	} else {
		res, err = rootNode.Query(selectorFlag)
		if err != nil {
			return fmt.Errorf("could not query node: %w", err)
		}
	}

	fmt.Printf("%v\n", res.Value)

	return nil
}

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag string

	cmd := &cobra.Command{
		Use:   "select -f <file> -p <json,yaml> -s <selector>",
		Short: "Select properties from the given file.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSelectCommand(fileFlag, parserFlag, selectorFlag)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")

	for _, f := range []string{"selector"} {
		if err := cmd.MarkFlagRequired(f); err != nil {
			panic("could not mark flag as required: " + f)
		}
	}

	return cmd
}
