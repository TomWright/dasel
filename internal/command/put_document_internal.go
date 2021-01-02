package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/storage"
	"io"
)

type putDocumentOpts struct {
	File           string
	Out            string
	ReadParser     string
	WriteParser    string
	Parser         string
	Selector       string
	DocumentString string
	DocumentParser string
	Reader         io.Reader
	Writer         io.Writer
	Multi          bool
	Compact        bool
}

func runPutDocumentCommand(opts putDocumentOpts, cmd *cobra.Command) error {
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

	documentParser, err := getPutDocumentParser(readParser, opts.DocumentParser)
	if err != nil {
		return err
	}

	documentValue, err := documentParser.FromBytes([]byte(opts.DocumentString))
	if err != nil {
		return fmt.Errorf("could not parse document: %w", err)
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, documentValue); err != nil {
			return fmt.Errorf("could not put document multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, documentValue); err != nil {
			return fmt.Errorf("could not put document value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File)
	if err != nil {
		return err
	}

	writeOptions := make([]storage.ReadWriteOption, 0)

	if opts.Compact {
		writeOptions = append(writeOptions, storage.PrettyPrintOption(false))
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, cmd, writeOptions...); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}

func putDocumentCommand() *cobra.Command {
	var documentParserFlag string

	cmd := &cobra.Command{
		Use:   "document -f <file> -d <document-parser> -s <selector> <document>",
		Short: "Put an entire document into the given document.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := putDocumentOpts{
				File:           cmd.Flag("file").Value.String(),
				Out:            cmd.Flag("out").Value.String(),
				ReadParser:     cmd.Flag("read").Value.String(),
				WriteParser:    cmd.Flag("write").Value.String(),
				Parser:         cmd.Flag("parser").Value.String(),
				Selector:       cmd.Flag("selector").Value.String(),
				DocumentParser: documentParserFlag,
				DocumentString: args[0],
			}
			opts.Multi, _ = cmd.Flags().GetBool("multiple")

			if opts.Selector == "" && len(args) > 1 {
				opts.Selector = args[0]
				opts.DocumentString = args[1]
			}
			return runPutDocumentCommand(opts, cmd)
		},
	}

	cmd.Flags().StringVarP(&documentParserFlag, "document-parser", "d", "", "The parser to use when reading the document")

	return cmd
}

func getPutDocumentParser(readParser storage.ReadParser, documentParserFlag string) (storage.ReadParser, error) {
	if documentParserFlag == "" {
		return readParser, nil
	}

	parser, err := storage.NewReadParserFromString(documentParserFlag)
	if err != nil {
		return nil, fmt.Errorf("could not get document parser: %w", err)
	}
	return parser, nil
}
