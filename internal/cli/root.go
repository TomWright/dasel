package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/internal"
	"github.com/tomwright/dasel/v3/parsing"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dasel",
		Short:   "Query and modify data structures using selectors",
		Long:    `dasel is a command-line utility to query and modify data structures using selectors.`,
		Version: internal.Version,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			selectorStr := ""
			if len(args) > 0 {
				selectorStr = args[0]
			}

			readerStr, _ := cmd.Flags().GetString("input")
			writerStr, _ := cmd.Flags().GetString("output")
			if writerStr == "" {
				writerStr = readerStr
			}

			reader, err := parsing.NewReader(parsing.Format(readerStr))
			writer, err := parsing.NewWriter(parsing.Format(writerStr))

			inputBytes, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}

			inputData, err := reader.Read(inputBytes)
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}

			outputData, err := execution.ExecuteSelector(selectorStr, inputData)
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

	cmd.AddCommand(manCommand(cmd))

	return cmd
}
