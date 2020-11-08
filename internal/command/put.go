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

type getRootNodeOpts struct {
	File   string
	Reader io.Reader
	Parser storage.Parser
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
	Parser storage.Parser
	File   string
	Out    string
	Writer io.Writer
	Plain  bool
}

func writeNodeToOutput(opts writeNodeToOutputOpts, cmd *cobra.Command) error {
	if opts.Writer == nil {
		switch {
		case opts.Out == "" && shouldReadFromStdin(opts.File):
			// No out flag and we read from stdin.
			opts.Writer = cmd.OutOrStdout()

		case opts.Out == "stdout", opts.Out == "-":
			// Out flag wants to write to stdout.
			opts.Writer = cmd.OutOrStdout()

		case opts.Out == "":
			// No out flag... write to the file we read from.
			f, err := os.Create(opts.File)
			if err != nil {
				return fmt.Errorf("could not open output file: %w", err)
			}
			defer f.Close()
			opts.Writer = f

		case opts.Out != "":
			// Out flag was set.
			f, err := os.Create(opts.Out)
			if err != nil {
				return fmt.Errorf("could not open output file: %w", err)
			}
			defer f.Close()
			opts.Writer = f
		}
	}

	if opts.Plain {
		if _, err := fmt.Fprintf(opts.Writer, "%v\n", opts.Node.InterfaceValue()); err != nil {
			return fmt.Errorf("could not write output: %w", err)
		}
		return nil
	}
	if err := storage.Write(opts.Parser, opts.Node.InterfaceValue(), opts.Node.OriginalValue, opts.Writer); err != nil {
		return fmt.Errorf("could not write to output file: %w", err)
	}

	return nil
}

type writeNodesToOutputOpts struct {
	Nodes  []*dasel.Node
	Parser storage.Parser
	File   string
	Out    string
	Writer io.Writer
	Plain  bool
}

func writeNodesToOutput(opts writeNodesToOutputOpts, cmd *cobra.Command) error {
	if opts.Writer == nil {
		switch {
		case opts.Out == "" && shouldReadFromStdin(opts.File):
			// No out flag and we read from stdin.
			opts.Writer = cmd.OutOrStdout()

		case opts.Out == "stdout", opts.Out == "-":
			// Out flag wants to write to stdout.
			opts.Writer = cmd.OutOrStdout()

		case opts.Out == "":
			// No out flag... write to the file we read from.
			f, err := os.Create(opts.File)
			if err != nil {
				return fmt.Errorf("could not open output file: %w", err)
			}
			defer f.Close()
			opts.Writer = f

		case opts.Out != "":
			// Out flag was set.
			f, err := os.Create(opts.Out)
			if err != nil {
				return fmt.Errorf("could not open output file: %w", err)
			}
			defer f.Close()
			opts.Writer = f
		}
	}

	buf := new(bytes.Buffer)

	for i, n := range opts.Nodes {
		subOpts := writeNodeToOutputOpts{
			Node:   n,
			Parser: opts.Parser,
			Writer: buf,
			Plain:  opts.Plain,
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

func putCommand() *cobra.Command {
	var fileFlag, selectorFlag, parserFlag, outFlag string
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
	cmd.PersistentFlags().StringVarP(&parserFlag, "parser", "p", "", "The parser to use with the given file.")
	cmd.PersistentFlags().StringVarP(&outFlag, "out", "o", "", "Output destination.")
	cmd.PersistentFlags().BoolVarP(&multiFlag, "multiple", "m", false, "Select multiple results.")

	_ = cmd.MarkPersistentFlagFilename("file")

	return cmd
}

type genericPutOptions struct {
	File      string
	Out       string
	Parser    string
	Selector  string
	Value     string
	ValueType string
	Init      func(genericPutOptions) genericPutOptions
	Reader    io.Reader
	Writer    io.Writer
	Multi     bool
}

func getGenericInit(cmd *cobra.Command, args []string) func(options genericPutOptions) genericPutOptions {
	return func(opts genericPutOptions) genericPutOptions {
		opts.File = cmd.Flag("file").Value.String()
		opts.Out = cmd.Flag("out").Value.String()
		opts.Parser = cmd.Flag("parser").Value.String()
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

	if err := writeNodeToOutput(writeNodeToOutputOpts{
		Node:   rootNode,
		Parser: parser,
		File:   opts.File,
		Out:    opts.Out,
		Writer: opts.Writer,
	}, cmd); err != nil {
		return fmt.Errorf("could not write output: %w", err)
	}

	return nil
}
