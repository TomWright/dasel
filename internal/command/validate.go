package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"sync"
)

type validationFile struct {
	File   string
	Parser string
}

type validationFileResult struct {
	File  validationFile
	Pass  bool
	Error error
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

			pass, err := validateFile(file)
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

func validateFile(file validationFile) (bool, error) {
	readParser, err := getReadParser(file.File, file.Parser, "")
	if err != nil {
		return false, err
	}
	_, err = getRootNode(getRootNodeOpts{
		File:   file.File,
		Parser: readParser,
		Reader: nil,
	}, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}

func validateCommand() *cobra.Command {
	var includeErrorFlag bool

	cmd := &cobra.Command{
		Use:   "validate <file> <file> <file>",
		Short: "Validate a list of files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			files := make([]validationFile, 0)
			for _, a := range args {
				files = append(files, validationFile{
					File:   a,
					Parser: "",
				})
			}

			return runValidateCommand(validateOptions{
				Files:        files,
				IncludeError: includeErrorFlag,
			}, cmd)
		},
	}

	cmd.Flags().BoolVar(&includeErrorFlag, "include-error", true, "Show error/validation information")

	_ = cmd.MarkFlagFilename("file")

	return cmd
}
