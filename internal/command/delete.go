package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
	"io"
)

type deleteOptions struct {
	File                string
	Parser              string
	ReadParser          string
	WriteParser         string
	Selector            string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
}

func runDeleteMultiCommand(cmd *cobra.Command, rootNode *dasel.Node, opts deleteOptions, writeParser storage.WriteParser, writeOptions []storage.ReadWriteOption) error {
	err := rootNode.DeleteMultiple(opts.Selector)

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		Cmd:      cmd,
		NullFlag: false,
	})
	if err != nil {
		return fmt.Errorf("could not delete multiple node: %w", err)
	}
	if written {
		return nil
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		Writer: opts.Writer,
	}, cmd, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}
	return nil
}

func runDeleteCommand(opts deleteOptions, cmd *cobra.Command) error {
	readParser, err := getReadParser(opts.File, opts.ReadParser, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:                opts.File,
		Parser:              readParser,
		Reader:              opts.Reader,
		MergeInputDocuments: opts.MergeInputDocuments,
	}, cmd)
	if err != nil {
		return err
	}

	if !rootNode.Value.IsValid() {
		rootNode = dasel.New(&storage.BasicSingleDocument{
			Value: map[string]interface{}{},
		})
	}

	if opts.Writer == nil {
		opts.Writer = cmd.OutOrStdout()
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, "-", opts.File)
	if err != nil {
		return err
	}

	writeOptions := make([]storage.ReadWriteOption, 0)

	if opts.Compact {
		writeOptions = append(writeOptions, storage.PrettyPrintOption(false))
	}

	if opts.Multi {
		return runDeleteMultiCommand(cmd, rootNode, opts, writeParser, writeOptions)
	}

	if err := rootNode.Delete(opts.Selector); err != nil {
		err = fmt.Errorf("could not delete node: %w", err)
	}

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		Cmd:      cmd,
		NullFlag: false,
	})
	if err != nil {
		return err
	}
	if written {
		return nil
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		Writer: opts.Writer,
	}, cmd, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func deleteCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, readParserFlag, writeParserFlag string
	var plainFlag, multiFlag, compactFlag, mergeInputDocumentsFlag bool

	cmd := &cobra.Command{
		Use:   "delete -f <file> -p <json,yaml> -s <selector>",
		Short: "Delete properties from the given file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if selectorFlag == "" && len(args) > 0 {
				selectorFlag = args[0]
				args = args[1:]
			}
			if plainFlag {
				writeParserFlag = "-"
			}
			return runDeleteCommand(deleteOptions{
				File:                fileFlag,
				Parser:              parserFlag,
				ReadParser:          readParserFlag,
				WriteParser:         writeParserFlag,
				Selector:            selectorFlag,
				Multi:               multiFlag,
				Compact:             compactFlag,
				MergeInputDocuments: mergeInputDocumentsFlag,
			}, cmd)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to delete from.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when deleting from the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "Shorthand for -r FORMAT -w FORMAT.")
	cmd.Flags().StringVarP(&readParserFlag, "read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringVarP(&writeParserFlag, "write", "w", "", "The parser to use when writing.")
	cmd.Flags().BoolVar(&plainFlag, "plain", false, "Alias of -w plain")
	cmd.Flags().BoolVarP(&multiFlag, "multiple", "m", false, "Delete multiple results.")
	cmd.Flags().BoolVar(&mergeInputDocumentsFlag, "merge-input-documents", false, "Merge multiple input documents into an array.")
	cmd.Flags().BoolVarP(&compactFlag, "compact", "c", false, "Compact the output by removing all pretty-printing where possible.")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
