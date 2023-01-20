package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
	"strconv"
)

func putCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put -t <string,int,float,bool,json,yaml,toml,xml,csv> -v <value> -f <file> -r <json,yaml,toml,xml,csv> <selector>",
		Short: "Write properties to the given file.",
		RunE:  putRunE,
		Args:  cobra.MaximumNArgs(1),
	}

	putFlags(cmd)

	return cmd
}

func putFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringP("read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringP("file", "f", "", "The file to query.")
	cmd.Flags().StringP("write", "w", "", "The parser to use when writing. Defaults to the read parser if not provided.")
	cmd.Flags().StringP("type", "t", "string", "The type of variable being written.")
	cmd.Flags().StringP("value", "v", "", "The value to write.")
	cmd.Flags().StringP("out", "o", "", "The file to write output to.")
	cmd.Flags().Bool("pretty", true, "Pretty print the output.")
	cmd.Flags().Bool("colour", false, "Print colourised output.")
	cmd.Flags().Bool("escape-html", false, "Escape HTML tags when writing output.")

	_ = cmd.MarkFlagFilename("file")
}

func putRunE(cmd *cobra.Command, args []string) error {
	selectorFlag, _ := cmd.Flags().GetString("selector")
	readParserFlag, _ := cmd.Flags().GetString("read")
	fileFlag, _ := cmd.Flags().GetString("file")
	writeParserFlag, _ := cmd.Flags().GetString("write")
	typeFlag, _ := cmd.Flags().GetString("type")
	valueFlag, _ := cmd.Flags().GetString("value")
	prettyPrintFlag, _ := cmd.Flags().GetBool("pretty")
	colourFlag, _ := cmd.Flags().GetBool("colour")
	escapeHTMLFlag, _ := cmd.Flags().GetBool("escape-html")
	outFlag, _ := cmd.Flags().GetString("out")

	opts := &putOptions{
		Read: &readOptions{
			Reader:   nil,
			Parser:   readParserFlag,
			FilePath: fileFlag,
		},
		Write: &writeOptions{
			Writer:      nil,
			Parser:      writeParserFlag,
			FilePath:    outFlag,
			PrettyPrint: prettyPrintFlag,
			Colourise:   colourFlag,
			EscapeHTML:  escapeHTMLFlag,
		},
		Selector:  selectorFlag,
		ValueType: typeFlag,
		Value:     valueFlag,
	}

	if opts.Selector == "" && len(args) > 0 {
		opts.Selector = args[0]
		args = args[1:]
	}

	if opts.Selector == "" {
		opts.Selector = "."
	}

	if opts.Write.FilePath == "" {
		opts.Write.FilePath = opts.Read.FilePath
	}

	return runPutCommand(opts, cmd)
}

type putOptions struct {
	Read      *readOptions
	Write     *writeOptions
	Selector  string
	ValueType string
	Value     string
}

func runPutCommand(opts *putOptions, cmd *cobra.Command) error {

	rootValue, err := opts.Read.rootValue(cmd)
	if err != nil {
		return err
	}

	var toSet interface{}

	switch opts.ValueType {
	case "string":
		toSet = opts.Value
	case "int":
		intValue, err := strconv.Atoi(opts.Value)
		if err != nil {
			return fmt.Errorf("invalid int value: %w", err)
		}
		toSet = intValue
	case "float":
		floatValue, err := strconv.ParseFloat(opts.Value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value: %w", err)
		}
		toSet = floatValue
	case "bool":
		toSet = dasel.ValueOf(dasel.IsTruthy(opts.Value))
	default:
		readParser, err := storage.NewReadParserFromString(opts.ValueType)
		if err != nil {
			return fmt.Errorf("unhandled value type")
		}
		docValue, err := readParser.FromBytes([]byte(opts.Value))
		if err != nil {
			return fmt.Errorf("could not parse document: %w", err)
		}
		toSet = docValue
	}

	value, err := dasel.Put(rootValue, opts.Selector, toSet)
	if err != nil {
		return err
	}

	if err := opts.Write.writeValue(cmd, opts.Read, value); err != nil {
		return err
	}

	return nil
}
