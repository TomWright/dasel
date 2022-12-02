package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"path/filepath"
	"sync"
)

func validateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <file> <file> <file>",
		Short: "Validate a list of files.",
		RunE:  validateRunE,
		Args:  cobra.ArbitraryArgs,
	}

	validateFlags(cmd)

	return cmd
}

func validateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("include-error", true, "Show error/validation information")
}

func validateRunE(cmd *cobra.Command, args []string) error {
	includeErrorFlag, _ := cmd.Flags().GetBool("include-error")

	files := make([]validationFile, 0)
	for _, a := range args {
		matches, err := filepath.Glob(a)
		if err != nil {
			return err
		}

		for _, m := range matches {
			files = append(files, validationFile{
				File:   m,
				Parser: "",
			})
		}
	}

	return runValidateCommand(validateOptions{
		Files:        files,
		IncludeError: includeErrorFlag,
	}, cmd)
}

type validateOptions struct {
	Files        []validationFile
	Reader       io.Reader
	Writer       io.Writer
	IncludeError bool
}

func runValidateCommand(opts validateOptions, cmd *cobra.Command) error {
	fileCount := len(opts.Files)

	wg := &sync.WaitGroup{}
	wg.Add(fileCount)

	results := make([]validationFileResult, fileCount)

	for i, f := range opts.Files {
		index := i
		file := f
		go func() {
			defer wg.Done()

			pass, err := validateFile(cmd, file)
			results[index] = validationFileResult{
				File:  file,
				Pass:  pass,
				Error: err,
			}
		}()
	}

	wg.Wait()

	failureCount := 0
	for _, result := range results {
		if !result.Pass {
			failureCount++
		}
	}

	// Set up our output writer if one wasn't provided.
	if opts.Writer == nil {
		if failureCount > 0 {
			opts.Writer = cmd.OutOrStderr()
		} else {
			opts.Writer = cmd.OutOrStdout()
		}
	}

	for _, result := range results {
		outputString := ""

		if result.Pass {
			outputString += "pass"
		} else {
			outputString += "fail"
		}

		outputString += " " + result.File.File

		if opts.IncludeError && result.Error != nil {
			outputString += " " + result.Error.Error()
		}

		if _, err := fmt.Fprintln(opts.Writer, outputString); err != nil {
			return fmt.Errorf("could not write output: %w", err)
		}
	}

	if failureCount > 0 {
		return fmt.Errorf("%d files failed validation", failureCount)
	}
	return nil
}

type validationFile struct {
	File   string
	Parser string
}

type validationFileResult struct {
	File  validationFile
	Pass  bool
	Error error
}

func validateFile(cmd *cobra.Command, file validationFile) (bool, error) {
	opts := readOptions{
		Parser:   file.Parser,
		FilePath: file.File,
	}
	_, err := opts.rootValue(cmd)

	return err == nil, err
}
