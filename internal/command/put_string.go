package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
)

func putStringCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "string -f <file> -s <selector> <value>",
		Short: "Update a string property in the given file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileFlag := cmd.Flag("file").Value.String()
			parserFlag := cmd.Flag("parser").Value.String()
			selectorFlag := cmd.Flag("selector").Value.String()

			parser, err := getParser(fileFlag, parserFlag)
			if err != nil {
				return err
			}
			value, err := storage.LoadFromFile(fileFlag, parser)
			if err != nil {
				return fmt.Errorf("could not load file: %w", err)
			}
			rootNode := dasel.New(value)

			updateValue, err := parseValue(args[0], "string")
			if err != nil {
				return err
			}

			if err := rootNode.Put(selectorFlag, updateValue); err != nil {
				return fmt.Errorf("could not put value: %w", err)
			}

			if err := storage.WriteToFile(fileFlag, parser, rootNode.Value); err != nil {
				return fmt.Errorf("could not write file: %w", err)
			}

			fmt.Println("updated string")

			return nil
		},
	}

	return cmd
}
