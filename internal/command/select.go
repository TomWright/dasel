package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
	"io"
	"reflect"
)

type selectOptions struct {
	File                string
	Parser              string
	ReadParser          string
	WriteParser         string
	Selector            string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	NullValueNotFound   bool
	Compact             bool
	DisplayLength       bool
	MergeInputDocuments bool
	FormatTemplate      string
}

func outputNodeLength(writer io.Writer, nodes ...*dasel.Node) error {
	for _, n := range nodes {
		val := n.Value
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		length := 0
		switch val.Kind() {
		case reflect.Map, reflect.Slice, reflect.String:
			length = val.Len()
		default:
			length = len(fmt.Sprint(val.Interface()))
		}
		if _, err := fmt.Fprintf(writer, "%d\n", length); err != nil {
			return err
		}
	}
	return nil
}

func runSelectMultiCommand(cmd *cobra.Command, rootNode *dasel.Node, opts selectOptions, writeParser storage.WriteParser, writeOptions []storage.ReadWriteOption) error {
	var results []*dasel.Node
	var err error
	if opts.Selector == "." {
		results = []*dasel.Node{rootNode}
	} else {
		results, err = rootNode.QueryMultiple(opts.Selector)
	}

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		Cmd:      cmd,
		NullFlag: opts.NullValueNotFound,
	})
	if err != nil {
		return fmt.Errorf("could not query multiple node: %w", err)
	}
	if written {
		return nil
	}

	if opts.DisplayLength {
		return outputNodeLength(opts.Writer, results...)
	}

	if err := writeNodesToOutput(writeNodesToOutputOpts{
		Nodes:          results,
		Parser:         writeParser,
		Writer:         opts.Writer,
		FormatTemplate: opts.FormatTemplate,
	}, cmd, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}
	return nil
}

func runSelectCommand(opts selectOptions, cmd *cobra.Command) error {
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

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, "-", opts.File, opts.FormatTemplate)
	if err != nil {
		return err
	}

	writeOptions := make([]storage.ReadWriteOption, 0)

	if opts.Compact {
		writeOptions = append(writeOptions, storage.PrettyPrintOption(false))
	}

	if opts.Multi {
		return runSelectMultiCommand(cmd, rootNode, opts, writeParser, writeOptions)
	}

	var res *dasel.Node

	if opts.Selector == "." {
		res = rootNode
	} else {
		res, err = rootNode.Query(opts.Selector)
		if err != nil {
			err = fmt.Errorf("could not query node: %w", err)
		}
	}

	written, err := customErrorHandling(customErrorHandlingOpts{
		File:     opts.File,
		Writer:   opts.Writer,
		Err:      err,
		Cmd:      cmd,
		NullFlag: opts.NullValueNotFound,
	})
	if err != nil {
		return err
	}
	if written {
		return nil
	}

	if opts.DisplayLength {
		return outputNodeLength(opts.Writer, res)
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:           res,
		Parser:         writeParser,
		Writer:         opts.Writer,
		FormatTemplate: opts.FormatTemplate,
	}, cmd, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, readParserFlag, writeParserFlag, formatTemplateFlag string
	var plainFlag, multiFlag, nullValueNotFoundFlag, compactFlag, lengthFlag, mergeInputDocumentsFlag bool

	cmd := &cobra.Command{
		Use:   "select -f <file> -p <json,yaml> -s <selector>",
		Short: "Select properties from the given file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if selectorFlag == "" && len(args) > 0 {
				selectorFlag = args[0]
				args = args[1:]
			}
			if plainFlag {
				writeParserFlag = "-"
			}
			return runSelectCommand(selectOptions{
				File:                fileFlag,
				Parser:              parserFlag,
				ReadParser:          readParserFlag,
				WriteParser:         writeParserFlag,
				Selector:            selectorFlag,
				Multi:               multiFlag,
				NullValueNotFound:   nullValueNotFoundFlag,
				Compact:             compactFlag,
				DisplayLength:       lengthFlag,
				MergeInputDocuments: mergeInputDocumentsFlag,
				FormatTemplate:      formatTemplateFlag,
			}, cmd)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "Shorthand for -r FORMAT -w FORMAT.")
	cmd.Flags().StringVarP(&readParserFlag, "read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringVarP(&writeParserFlag, "write", "w", "", "The parser to use when writing.")
	cmd.Flags().BoolVar(&plainFlag, "plain", false, "Alias of -w plain")
	cmd.Flags().BoolVarP(&multiFlag, "multiple", "m", false, "Select multiple results.")
	cmd.Flags().BoolVarP(&nullValueNotFoundFlag, "null", "n", false, "Output null instead of value not found errors.")
	cmd.Flags().BoolVar(&lengthFlag, "length", false, "Output the length of the selected value.")
	cmd.Flags().BoolVar(&mergeInputDocumentsFlag, "merge-input-documents", false, "Merge multiple input documents into an array.")
	cmd.Flags().BoolVarP(&compactFlag, "compact", "c", false, "Compact the output by removing all pretty-printing where possible.")
	cmd.Flags().StringVar(&formatTemplateFlag, "format", "", "Formatting template to use when writing results.")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
