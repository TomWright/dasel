package command

import (
	"fmt"
	"github.com/spf13/cobra"
	dasel "github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
	"strconv"
	"strings"
)

func parseValue(value string, valueType string) (interface{}, error) {
	switch strings.ToLower(valueType) {
	case "string", "str":
		return value, nil
	case "int", "integer":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse int [%s]: %w", value, err)
		}
		return val, nil
	case "bool", "boolean":
		switch strings.ToLower(value) {
		case "true", "t", "yes", "y", "1":
			return true, nil
		case "false", "f", "no", "n", "0":
			return false, nil
		default:
			return nil, fmt.Errorf("could not parse bool [%s]: unhandled value", value)
		}
	default:
		return nil, fmt.Errorf("unhandled type: %s", valueType)
	}
}

func putCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, typeFlag string

	cmd := &cobra.Command{
		Use:   "put -f <file> -s <selector> -t <string|int|bool> <value>",
		Short: "Update properties in the given file.",
		Args:  cobra.ExactArgs(1),
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

			updateValue, err := parseValue(args[0], typeFlag)
			if err != nil {
				return fmt.Errorf("invalid value: %w", err)
			}

			if err := rootNode.Put(selectorFlag, updateValue); err != nil {
				return fmt.Errorf("could not put value: %w", err)
			}

			fmt.Println("updated")

			return nil
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")
	cmd.Flags().StringVarP(&parserFlag, "type", "t", "string", "The type of variable we are updating.")

	for _, f := range []string{"file", "selector", "type"} {
		if err := cmd.MarkFlagRequired(f); err != nil {
			panic("could not mark flag as required: " + f)
		}
	}

	return cmd
}
