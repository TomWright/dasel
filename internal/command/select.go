package command

import (
	"fmt"
	"github.com/spf13/cobra"
	dasel "github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
)

func selectCommand() *cobra.Command {
	var file, selector, parser string

	cmd := &cobra.Command{
		Use:   "select -f <file> -s <selector>",
		Short: "Select properties from the given files.",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			parser, err := storage.FromString(parser)
			if err != nil {
				return err
			}
			value, err := parser.LoadFromFile(file)
			if err != nil {
				return fmt.Errorf("could not load file: %w", err)
			}
			rootNode := dasel.New(value)
			res, err := rootNode.Query(selector)
			if err != nil {
				return fmt.Errorf("could not query node: %w", err)
			}

			fmt.Printf("%v\n", res.Value)

			return nil
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selector, "selector", "s", "", "The selector to use when querying.")
	cmd.Flags().StringVarP(&parser, "parser", "p", "yaml", "The parser to use with the given file.")

	return cmd
}
