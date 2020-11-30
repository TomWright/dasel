package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"io"
)

type selectOptions struct {
	File        string
	Parser      string
	ReadParser  string
	WriteParser string
	Selector    string
	Reader      io.Reader
	Writer      io.Writer
	Multi       bool
}

func runSelectCommand(opts selectOptions, cmd *cobra.Command) error {
	readParser, err := getReadParser(opts.File, opts.ReadParser, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:   opts.File,
		Parser: readParser,
		Reader: opts.Reader,
	}, cmd)
	if err != nil {
		return err
	}

	if opts.Writer == nil {
		opts.Writer = cmd.OutOrStdout()
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, "-", opts.File)
	if err != nil {
		return err
	}

	if opts.Multi {
		var results []*dasel.Node
		if opts.Selector == "." {
			results = []*dasel.Node{rootNode}
		} else {
			results, err = rootNode.QueryMultiple(opts.Selector)
			if err != nil {
				return fmt.Errorf("could not query multiple node: %w", err)
			}
		}

		if err := writeNodesToOutput(writeNodesToOutputOpts{
			Nodes:  results,
			Parser: writeParser,
			Writer: opts.Writer,
		}, cmd); err != nil {
			return fmt.Errorf("could not write output: %w", err)
		}
		return nil
	}

	var res *dasel.Node

	if opts.Selector == "." {
		res = rootNode
	} else {
		res, err = rootNode.Query(opts.Selector)
		if err != nil {
			return fmt.Errorf("could not query node: %w", err)
		}
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   res,
		Parser: writeParser,
		Writer: opts.Writer,
	}, cmd); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, readParserFlag, writeParserFlag string
	var plainFlag, multiFlag bool

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
				File:        fileFlag,
				Parser:      parserFlag,
				ReadParser:  readParserFlag,
				WriteParser: writeParserFlag,
				Selector:    selectorFlag,
				Multi:       multiFlag,
			}, cmd)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "Shorthand for `-r FORMAT -w FORMAT`.")
	cmd.Flags().StringVarP(&readParserFlag, "read", "r", "", "The parser to use when reading.")
	cmd.Flags().StringVarP(&writeParserFlag, "write", "w", "", "The parser to use when writing.")
	cmd.Flags().BoolVar(&plainFlag, "plain", false, "Do not format output to the output data format.")
	cmd.Flags().BoolVarP(&multiFlag, "multiple", "m", false, "Select multiple results.")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
