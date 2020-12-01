package command

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel"
	"github.com/tomwright/dasel/internal/storage"
	"io"
	"os"
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
	return fileFlag == "" || fileFlag == "stdin" || fileFlag == "-"
}

func shouldWriteToStdout(fileFlag string, outFlag string) bool {
	return (outFlag == "stdout" || outFlag == "-") || outFlag == "" && shouldReadFromStdin(fileFlag)
}

func getReadParser(fileFlag string, readParserFlag string, parserFlag string) (storage.ReadParser, error) {
	useStdin := shouldReadFromStdin(fileFlag)

	if readParserFlag == "" {
		readParserFlag = parserFlag
	}

	if useStdin && readParserFlag == "" {
		return nil, fmt.Errorf("read parser flag required when reading from stdin")
	}

	if readParserFlag == "" {
		parser, err := storage.NewReadParserFromFilename(fileFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get read parser from filename: %w", err)
		}
		return parser, nil
	}
	parser, err := storage.NewReadParserFromString(readParserFlag)
	if err != nil {
		return nil, fmt.Errorf("could not get read parser: %w", err)
	}
	return parser, nil
}

func getWriteParser(readParser storage.ReadParser, writeParserFlag string, parserFlag string, outFlag string, fileFlag string) (storage.WriteParser, error) {
	if writeParserFlag == "" {
		writeParserFlag = parserFlag
	}

	if writeParserFlag != "" {
		parser, err := storage.NewWriteParserFromString(writeParserFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser: %w", err)
		}
		return parser, nil
	}

	if !shouldWriteToStdout(fileFlag, outFlag) {
		p, err := storage.NewWriteParserFromFilename(outFlag)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser from filename: %w", err)
		}
		return p, nil
	}

	if p, ok := readParser.(storage.WriteParser); ok {
		return p, nil
	}
	return nil, fmt.Errorf("read parser cannot be used to write. please specify a write parser")
}

type getRootNodeOpts struct {
	File   string
	Reader io.Reader
	Parser storage.ReadParser
}

func getRootNode(opts getRootNodeOpts, cmd *cobra.Command) (*dasel.Node, error) {
	if opts.Reader == nil {
		if shouldReadFromStdin(opts.File) {
			opts.Reader = cmd.InOrStdin()
		} else {
			f, err := os.Open(opts.File)
			if err != nil {
				return nil, fmt.Errorf("could not open input file: %w", err)
			}
			defer f.Close()
			opts.Reader = f
		}
	}

	value, err := storage.Load(opts.Parser, opts.Reader)
	if err != nil {
		return nil, fmt.Errorf("could not load input: %w", err)
	}

	return dasel.New(value), nil
}

type writeNodeToOutputOpts struct {
	Node   *dasel.Node
	Parser storage.WriteParser
	File   string
	Out    string
	Writer io.Writer
}

func writeNodeToOutput(opts writeNodeToOutputOpts, cmd *cobra.Command) error {
	writer, writerCleanUp, err := getOutputWriter(cmd, opts.Writer, opts.File, opts.Out)
	if err != nil {
		return err
	}
	opts.Writer = writer
	defer writerCleanUp()

	if err := storage.Write(opts.Parser, opts.Node.InterfaceValue(), opts.Node.OriginalValue, opts.Writer); err != nil {
		return fmt.Errorf("could not write to output file: %w", err)
	}

	return nil
}

type writeNodesToOutputOpts struct {
	Nodes  []*dasel.Node
	Parser storage.WriteParser
	File   string
	Out    string
	Writer io.Writer
}

func writeNodesToOutput(opts writeNodesToOutputOpts, cmd *cobra.Command) error {
	writer, writerCleanUp, err := getOutputWriter(cmd, opts.Writer, opts.File, opts.Out)
	if err != nil {
		return err
	}
	opts.Writer = writer
	defer writerCleanUp()

	buf := new(bytes.Buffer)

	for i, n := range opts.Nodes {
		subOpts := writeNodeToOutputOpts{
			Node:   n,
			Parser: opts.Parser,
			Writer: buf,
		}
		if err := writeNodeToOutput(subOpts, cmd); err != nil {
			return fmt.Errorf("could not write node %d to output: %w", i, err)
		}
	}

	if _, err := io.Copy(opts.Writer, buf); err != nil {
		return fmt.Errorf("could not copy buffer to real output: %w", err)
	}

	return nil
}

func getOutputWriter(cmd *cobra.Command, in io.Writer, file string, out string) (io.Writer, func(), error) {
	if in == nil {
		switch {
		case shouldWriteToStdout(file, out):
			return cmd.OutOrStdout(), func() {}, nil

		case out == "":
			// No out flag... write to the file we read from.
			f, err := os.Create(file)
			if err != nil {
				return nil, nil, fmt.Errorf("could not open output file: %w", err)
			}
			return f, func() {
				_ = f.Close()
			}, nil

		case out != "":
			// Out flag was set.
			f, err := os.Create(out)
			if err != nil {
				return nil, nil, fmt.Errorf("could not open output file: %w", err)
			}
			return f, func() {
				_ = f.Close()
			}, nil
		}
	}
	return in, func() {}, nil
}

func putCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, readParserFlag, writeParserFlag, outFlag string
	var multiFlag bool

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
	cmd.PersistentFlags().StringVarP(&parserFlag, "parser", "p", "", "Shorthand for -r FORMAT -w FORMAT.")
	cmd.PersistentFlags().StringVarP(&readParserFlag, "read", "r", "", "The parser to use when reading.")
	cmd.PersistentFlags().StringVarP(&writeParserFlag, "write", "w", "", "The parser to use when writing.")
	cmd.PersistentFlags().StringVarP(&outFlag, "out", "o", "", "Output destination.")
	cmd.PersistentFlags().BoolVarP(&multiFlag, "multiple", "m", false, "Select multiple results.")

	_ = cmd.MarkPersistentFlagFilename("file")

	return cmd
}

type genericPutOptions struct {
	File        string
	Out         string
	Parser      string
	ReadParser  string
	WriteParser string
	Selector    string
	Value       string
	ValueType   string
	Init        func(genericPutOptions) genericPutOptions
	Reader      io.Reader
	Writer      io.Writer
	Multi       bool
}

func getGenericInit(cmd *cobra.Command, args []string) func(options genericPutOptions) genericPutOptions {
	return func(opts genericPutOptions) genericPutOptions {
		opts.File = cmd.Flag("file").Value.String()
		opts.Out = cmd.Flag("out").Value.String()
		opts.Parser = cmd.Flag("parser").Value.String()
		opts.ReadParser = cmd.Flag("read").Value.String()
		opts.WriteParser = cmd.Flag("write").Value.String()
		opts.Selector = cmd.Flag("selector").Value.String()
		opts.Multi, _ = cmd.Flags().GetBool("multiple")

		if opts.Selector == "" && len(args) > 0 {
			opts.Selector = args[0]
			args = args[1:]
		}

		if len(args) > 0 {
			opts.Value = args[0]
		}

		return opts
	}
}

func runGenericPutCommand(opts genericPutOptions, cmd *cobra.Command) error {
	if opts.Init != nil {
		opts = opts.Init(opts)
	}
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

	updateValue, err := parseValue(opts.Value, opts.ValueType)
	if err != nil {
		return err
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File)
	if err != nil {
		return err
	}

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: writeParser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, cmd); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}
