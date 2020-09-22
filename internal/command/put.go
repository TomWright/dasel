package command

import (
	"fmt"
	"github.com/spf13/cobra"
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

func getParser(fileFlag string, parserFlag string) (storage.Parser, error) {
	if parserFlag == "" {
		parser, err := storage.NewParserFromFilename(fileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get parser from filename: %w", err)
		}
		return parser, nil
	} else {
		parser, err := storage.NewParserFromString(parserFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get parser: %w", err)
		}
		return parser, nil
	}
}

func putCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag string

	cmd := &cobra.Command{
		Use:   "put -f <file> -s <selector>",
		Short: "Update properties in the given file.",
	}

	cmd.AddCommand(
		putStringCommand(),
		putBoolCommand(),
		putIntCommand(),
		putObjectCommand(),
	)

	cmd.PersistentFlags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.PersistentFlags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.PersistentFlags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")

	for _, f := range []string{"file", "selector"} {
		if err := cmd.MarkPersistentFlagRequired(f); err != nil {
			panic("could not mark flag as required: " + f)
		}
	}

	return cmd
}
