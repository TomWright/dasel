package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
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

func shouldReadFromStdin(fileFlag string) bool {
	return fileFlag == ""
}

func getParser(fileFlag string, parserFlag string) (storage.Parser, error) {
	useStdin := shouldReadFromStdin(fileFlag)
	if useStdin && parserFlag == "" {
		return nil, fmt.Errorf("parser flag required when reading from stdin")
	}

	if parserFlag == "" {
		parser, err := storage.NewParserFromFilename(fileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get parser from filename: %w", err)
		}
		return parser, nil
	}
	parser, err := storage.NewParserFromString(parserFlag)
	if err != nil {
		return nil, fmt.Errorf("could not get parser: %w", err)
	}
	return parser, nil
}

func getRootNode(fileFlag string, parser storage.Parser) (*dasel.Node, error) {
	useStdin := shouldReadFromStdin(fileFlag)
	if useStdin {
		value, err := storage.LoadFromStdin(parser)
		if err != nil {
			return nil, fmt.Errorf("could not load file: %w", err)
		}
		return dasel.New(value), nil
	}
	value, err := storage.LoadFromFile(fileFlag, parser)
	if err != nil {
		return nil, fmt.Errorf("could not load file: %w", err)
	}
	return dasel.New(value), nil
}

func writeNodeToOutput(n *dasel.Node, parser storage.Parser, fileFlag string, outFlag string) error {
	if shouldReadFromStdin(fileFlag) && outFlag == "" {
		outFlag = "stdout"
	}
	switch outFlag {
	case "":
		return storage.WriteToFile(fileFlag, parser, n.Value)
	case "stdout":
		return storage.WriteToStdout(parser, n.Value)
	default:
		return storage.WriteToFile(outFlag, parser, n.Value)
	}
}

func putCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, outFlag string

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
	cmd.PersistentFlags().StringVarP(&outFlag, "out", "o", "", "Output destination.")

	for _, f := range []string{"selector"} {
		if err := cmd.MarkPersistentFlagRequired(f); err != nil {
			panic("could not mark flag as required: " + f)
		}
	}

	return cmd
}
