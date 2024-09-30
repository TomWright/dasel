package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v3"
)

func selectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dasel -f <file> -r <json,yaml,toml,xml,csv> <selector>",
		Short: "Select properties from the given file.",
		RunE:  selectRunE,
		Args:  cobra.MaximumNArgs(1),
	}

	selectFlags(cmd)

	return cmd
}

func selectFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringP("read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringP("file", "f", "", "The file to query.")
	cmd.Flags().StringP("write", "w", "", "The parser to use when writing. Defaults to the read parser if not provided.")
	cmd.Flags().Bool("pretty", true, "Pretty print the output.")
	cmd.Flags().Bool("colour", false, "Print colourised output.")
	cmd.Flags().Bool("escape-html", false, "Escape HTML tags when writing output.")
	cmd.Flags().Int("indent", 2, "The indention level when writing files.")
	cmd.Flags().String("csv-comma", ",", "Comma separator to use when working with csv files.")
	cmd.Flags().String("csv-write-comma", "", "Comma separator used when writing csv files. Overrides csv-comma when writing.")
	cmd.Flags().String("csv-comment", "", "Comma separator used when reading csv files.")
	cmd.Flags().Bool("csv-crlf", false, "True to use CRLF when writing CSV files.")

	_ = cmd.MarkFlagFilename("file")
}

func selectRunE(cmd *cobra.Command, args []string) error {
	selectorFlag, _ := cmd.Flags().GetString("selector")
	readParserFlag, _ := cmd.Flags().GetString("read")
	fileFlag, _ := cmd.Flags().GetString("file")
	writeParserFlag, _ := cmd.Flags().GetString("write")
	prettyPrintFlag, _ := cmd.Flags().GetBool("pretty")
	colourFlag, _ := cmd.Flags().GetBool("colour")
	escapeHTMLFlag, _ := cmd.Flags().GetBool("escape-html")
	indent, _ := cmd.Flags().GetInt("indent")
	csvComma, _ := cmd.Flags().GetString("csv-comma")
	csvWriteComma, _ := cmd.Flags().GetString("csv-write-comma")
	csvComment, _ := cmd.Flags().GetString("csv-comment")
	csvCRLF, _ := cmd.Flags().GetBool("csv-crlf")

	opts := &selectOptions{
		Read: &readOptions{
			Reader:     nil,
			Parser:     readParserFlag,
			FilePath:   fileFlag,
			CsvComma:   csvComma,
			CsvComment: csvComment,
		},
		Write: &writeOptions{
			Writer:      nil,
			Parser:      writeParserFlag,
			FilePath:    "stdout",
			PrettyPrint: prettyPrintFlag,
			Colourise:   colourFlag,
			EscapeHTML:  escapeHTMLFlag,
			Indent:      indent,
			CsvComma:    csvWriteComma,
			CsvUseCRLF:  csvCRLF,
		},
		Selector: selectorFlag,
	}

	if opts.Selector == "" && len(args) > 0 {
		opts.Selector = args[0]
		args = args[1:]
	}

	if opts.Selector == "" {
		opts.Selector = "."
	}

	return runSelectCommand(opts, cmd)
}

type selectOptions struct {
	Read     *readOptions
	Write    *writeOptions
	Selector string
}

func runSelectCommand(opts *selectOptions, cmd *cobra.Command) error {

	rootValue, err := opts.Read.rootValue(cmd)
	if err != nil {
		return err
	}

	values, err := dasel.Select(rootValue, opts.Selector)
	if err != nil {
		return err
	}

	if err := opts.Write.writeValues(cmd, opts.Read, values); err != nil {
		return err
	}

	return nil
}
