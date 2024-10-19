package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/internal"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

// ErrBadVarArg is returned when an invalid variable argument is provided.
var ErrBadVarArg = errors.New("invalid variable format, expect foo=bar, or foo=format:path")

// varOptFromArg attempts to parse a variable declaration from a commandline argument.
func varOptFromArg(arg string, r parsing.Reader) (execution.ExecuteOptionFn, error) {
	kv := strings.SplitN(arg, "=", 2)
	if len(kv) != 2 {
		return nil, ErrBadVarArg
	}

	varName := kv[0]

	formatData := strings.SplitN(kv[1], ":", 2)

	reader := r
	filepath := kv[1]

	if len(formatData) == 2 {
		var err error
		reader, err = parsing.NewReader(parsing.Format(formatData[0]))
		if err != nil {
			return nil, err
		}
		filepath = formatData[1]
	} else if reader == nil {
		return nil, fmt.Errorf("variable file format required")
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	inputBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	value, err := reader.Read(inputBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents for variable %q: %w", varName, err)
	}

	return execution.WithVariable(varName, value), nil
}

// RootCmd returns the root cli command.
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dasel [flags] [variableName=fileFormat:filePath] selector",
		Short:   "Query and modify data structures using selectors",
		Long:    `dasel is a command-line utility to query and modify data structures using selectors.`,
		Version: internal.Version,
		Args:    nil,
		Example: `dasel -o json foo=json:bar.json '{"x": $foo.x, "y": $foo.x + 1}'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			selectorStr := ""
			if len(args) > 0 {
				selectorStr = args[len(args)-1]
				args = args[0 : len(args)-1]
			}

			var opts []execution.ExecuteOptionFn

			readerStr, _ := cmd.Flags().GetString("input")
			writerStr, _ := cmd.Flags().GetString("output")
			if writerStr == "" {
				writerStr = readerStr
			}

			var reader parsing.Reader
			var err error
			if len(readerStr) > 0 {
				reader, err = parsing.NewReader(parsing.Format(readerStr))
				if err != nil {
					return fmt.Errorf("failed to get input reader: %w", err)
				}
			}

			writer, err := parsing.NewWriter(parsing.Format(writerStr))
			if err != nil {
				return fmt.Errorf("failed to get output writer: %w", err)
			}

			for _, a := range args {
				o, err := varOptFromArg(a, reader)
				if err != nil {
					return fmt.Errorf("failed to process variable: %w", err)
				}
				opts = append(opts, o)
			}

			var inputBytes []byte
			if reader != nil {
				inputBytes, err = io.ReadAll(cmd.InOrStdin())
				if err != nil {
					return fmt.Errorf("error reading input: %w", err)
				}
			}

			var inputData *model.Value
			if len(inputBytes) > 0 {
				inputData, err = reader.Read(inputBytes)
				if err != nil {
					return fmt.Errorf("error reading input: %w", err)
				}
				opts = append(opts, execution.WithVariable("root", inputData))
			}

			options := execution.NewOptions(opts...)

			outputData, err := execution.ExecuteSelector(selectorStr, inputData, options)
			if err != nil {
				return err
			}

			outputBytes, err := writer.Write(outputData)
			if err != nil {
				return fmt.Errorf("error writing output: %w", err)
			}

			_, err = cmd.OutOrStdout().Write(outputBytes)
			if err != nil {
				return fmt.Errorf("error writing output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP("input", "i", "", "The format of the input data. Can be one of: json, yaml, toml, xml, csv")
	cmd.Flags().StringP("output", "o", "", "The format of the output data. Can be one of: json, yaml, toml, xml, csv")

	// TODO : apply fallback to root cmd
	//cmd.AddCommand(manCommand(cmd))

	return cmd
}
