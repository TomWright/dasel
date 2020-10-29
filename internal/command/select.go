package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"io"
)

type selectOptions struct {
	File     string
	Parser   string
	Selector string
	Reader   io.Reader
	Writer   io.Writer
	Plain    bool
}

func runSelectCommand(opts selectOptions, cmd *cobra.Command) error {
	parser, err := getParser(opts.File, opts.Parser)
	if err != nil {
		return err
	}
	rootNode, err := getRootNode(getRootNodeOpts{
		File:   opts.File,
		Parser: parser,
		Reader: opts.Reader,
	}, cmd)
	if err != nil {
		return err
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

	if opts.Writer == nil {
		opts.Writer = cmd.OutOrStdout()
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   res,
		Parser: parser,
		Writer: opts.Writer,
		Plain:  opts.Plain,
	}, cmd); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func selectCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag string
	var plainFlag bool

	cmd := &cobra.Command{
		Use:   "select -f <file> -p <json,yaml> -s <selector>",
		Short: "Select properties from the given file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if selectorFlag == "" && len(args) > 0 {
				selectorFlag = args[0]
				args = args[1:]
			}
			return runSelectCommand(selectOptions{
				File:     fileFlag,
				Parser:   parserFlag,
				Selector: selectorFlag,
				Plain:    plainFlag,
			}, cmd)
		},
	}

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "The file to query.")
	cmd.Flags().StringVarP(&selectorFlag, "selector", "s", "", "The selector to use when querying the data structure.")
	cmd.Flags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")
	cmd.Flags().BoolVar(&plainFlag, "plain", false, "Do not format output to the output data format.")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
