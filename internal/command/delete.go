package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
)

func deleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete -f <file> -r <json,yaml,toml,xml,csv> <selector>",
		Short: "Delete properties from the given file.",
		RunE:  deleteRunE,
		Args:  cobra.MaximumNArgs(1),
	}

	deleteFlags(cmd)

	return cmd
}

func deleteFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringP("read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringP("file", "f", "", "The file to query.")
	cmd.Flags().StringP("write", "w", "", "The parser to use when writing. Defaults to the read parser if not provided.")
	cmd.Flags().Bool("pretty", true, "Pretty print the output.")
	cmd.Flags().Bool("colour", false, "Print colourised output.")
	cmd.Flags().Bool("escape-html", false, "Escape HTML tags when writing output.")

	_ = cmd.MarkFlagFilename("file")
}

func deleteRunE(cmd *cobra.Command, args []string) error {
	selectorFlag, _ := cmd.Flags().GetString("selector")
	readParserFlag, _ := cmd.Flags().GetString("read")
	fileFlag, _ := cmd.Flags().GetString("file")
	writeParserFlag, _ := cmd.Flags().GetString("write")
	prettyPrintFlag, _ := cmd.Flags().GetBool("pretty")
	colourFlag, _ := cmd.Flags().GetBool("colour")
	escapeHTMLFlag, _ := cmd.Flags().GetBool("escape-html")

	opts := &deleteOptions{
		Read: &readOptions{
			Reader:   nil,
			Parser:   readParserFlag,
			FilePath: fileFlag,
		},
		Write: &writeOptions{
			Writer:      nil,
			Parser:      writeParserFlag,
			FilePath:    "stdout",
			PrettyPrint: prettyPrintFlag,
			Colourise:   colourFlag,
			EscapeHTML:  escapeHTMLFlag,
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

	return deletePutCommand(opts, cmd)
}

type deleteOptions struct {
	Read     *readOptions
	Write    *writeOptions
	Selector string
}

func deletePutCommand(opts *deleteOptions, cmd *cobra.Command) error {

	rootValue, err := opts.Read.rootValue(cmd)
	if err != nil {
		return err
	}

	c := dasel.NewContext(rootValue, opts.Selector).WithCreateWhenMissing(true)

	values, err := c.Run()
	if err != nil {
		return err
	}

	for _, v := range values {
		v.Delete()
	}

	// There are issues deleting
	if err := opts.Write.writeValue(cmd, opts.Read, c.Data(dasel.WithoutDeletePlaceholders)); err != nil {
		return err
	}

	return nil
}
