package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/oflag"
	"github.com/tomwright/dasel/storage"
	"io"
	"strings"
)

type putObjectOpts struct {
	File                string
	Out                 string
	ReadParser          string
	WriteParser         string
	Parser              string
	Selector            string
	InputTypes          []string
	InputValues         []string
	Reader              io.Reader
	Writer              io.Writer
	Multi               bool
	Compact             bool
	MergeInputDocuments bool
	EscapeHTML          bool
}

func getMapFromTypesValues(inputTypes []string, inputValues []string) (map[string]interface{}, error) {
	if len(inputTypes) != len(inputValues) {
		return nil, fmt.Errorf("exactly %d types are required, got %d", len(inputValues), len(inputTypes))
	}

	updateValue := map[string]interface{}{}

	for k, arg := range inputValues {
		splitArg := strings.Split(arg, "=")
		name := splitArg[0]
		value := strings.Join(splitArg[1:], "=")
		parsedValue, err := parseValue(value, inputTypes[k])
		if err != nil {
			return nil, fmt.Errorf("could not parse value [%s]: %w", name, err)
		}
		updateValue[name] = parsedValue
	}

	return updateValue, nil
}

func runPutObjectCommand(opts putObjectOpts, cmd *cobra.Command) error {
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

	updateValue, err := getMapFromTypesValues(opts.InputTypes, opts.InputValues)
	if err != nil {
		return err
	}

	if opts.Multi {
		if err := rootNode.PutMultiple(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put object multi value: %w", err)
		}
	} else {
		if err := rootNode.Put(opts.Selector, updateValue); err != nil {
			return fmt.Errorf("could not put object value: %w", err)
		}
	}

	writeParser, err := getWriteParser(readParser, opts.WriteParser, opts.Parser, opts.Out, opts.File, "")
	if err != nil {
		return err
	}

	writeOptions := []storage.ReadWriteOption{
		storage.EscapeHTMLOption(opts.EscapeHTML),
	}

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

func putObjectCommand() *cobra.Command {
	typeList := oflag.NewStringList()

	cmd := &cobra.Command{
		Use:   "object -f <file> -s <selector> <value>",
		Short: "Put an object in the given document.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := putObjectOpts{
				File:        cmd.Flag("file").Value.String(),
				Out:         cmd.Flag("out").Value.String(),
				ReadParser:  cmd.Flag("read").Value.String(),
				WriteParser: cmd.Flag("write").Value.String(),
				Parser:      cmd.Flag("parser").Value.String(),
				Selector:    cmd.Flag("selector").Value.String(),
				InputTypes:  typeList.Strings,
				InputValues: args,
			}
			opts.Multi, _ = cmd.Flags().GetBool("multiple")
			opts.Compact, _ = cmd.Flags().GetBool("compact")
			opts.MergeInputDocuments, _ = cmd.Flags().GetBool("merge-input-documents")
			opts.EscapeHTML, _ = cmd.Flags().GetBool("escape-html")

			if opts.Selector == "" && len(opts.InputValues) > 0 {
				opts.Selector = opts.InputValues[0]
				opts.InputValues = opts.InputValues[1:]
			}
			return runPutObjectCommand(opts, cmd)
		},
	}

	cmd.Flags().VarP(typeList, "type", "t", "Types of the variables in the object.")

	return cmd
}
