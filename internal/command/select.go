package command

import (
	"fmt"
	"github.com/spf13/cobra"
	dasel "github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
)

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag string

	cmd := &cobra.Command{
		Use:   "select -f <file> -s <selector>",
		Short: "Select properties from the given file.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var parser storage.Parser
			var err error
			if parserFlag == "" {
				parser, err = storage.NewParserFromFilename(fileFlag)
				if err != nil {
					return fmt.Errorf("could not get parser from filename: %w", err)
				}
			} else {
				parser, err = storage.NewParserFromString(parserFlag)
				if err != nil {
					return fmt.Errorf("could not get parser: %w", err)
				}
			}
			value, err := storage.LoadFromFile(fileFlag, parser)
			if err != nil {
				return fmt.Errorf("could not load file: %w", err)
			}
			rootNode := dasel.New(value)
			res, err := rootNode.Query(selectorFlag)
			if err != nil {
				return fmt.Errorf("could not query node: %w", err)
			}

			fmt.Printf("%v\n", res.Value)

			return nil
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")

	for _, f := range []string{"file", "selector"} {
		if err := cmd.MarkFlagRequired(f); err != nil {
			panic("could not mark flag as required: " + f)
		}
	}

	return cmd
}
