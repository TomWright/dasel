package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/v2"
	"github.com/tomwright/dasel/v2/storage"
	"io"
	"os"
	"strconv"
)

type readOptions struct {
	// Reader is an io.Reader that we should read from instead of FilePath.
	Reader io.Reader
	// Parser is the name of the parser we should use when reading.
	Parser string
	// FilePath is the path to the source file.
	FilePath string
	// CsvComma is the comma character used when reading CSV files.
	CsvComma string
	// CsvComment is the comment character used when reading CSV files.
	CsvComment string
}

func (o *readOptions) readFromStdin() bool {
	return o.FilePath == "" || o.FilePath == "stdin" || o.FilePath == "-"
}

func (o *readOptions) readParser() (storage.ReadParser, error) {
	useStdin := o.readFromStdin()

	if useStdin && o.Parser == "" {
		return nil, fmt.Errorf("read parser required when reading from stdin")
	}

	if o.Parser == "" {
		parser, err := storage.NewReadParserFromFilename(o.FilePath)
		if err != nil {
			return nil, fmt.Errorf("could not get read parser from filename: %w", err)
		}
		return parser, nil
	}
	parser, err := storage.NewReadParserFromString(o.Parser)
	if err != nil {
		return nil, fmt.Errorf("could not get read parser: %w", err)
	}
	return parser, nil
}

func (o *readOptions) rootValue(cmd *cobra.Command) (dasel.Value, error) {
	parser, err := o.readParser()
	if err != nil {
		return dasel.Value{}, fmt.Errorf("could not get read parser: %w", err)
	}

	options := make([]storage.ReadWriteOption, 0)
	if o.CsvComma != "" {
		options = append(options, storage.CsvCommaOption([]rune(o.CsvComma)[0]))
	}
	if o.CsvComment != "" {
		options = append(options, storage.CsvCommentOption([]rune(o.CsvComment)[0]))
	}

	reader := o.Reader
	if reader == nil {
		if o.readFromStdin() {
			reader = cmd.InOrStdin()
		} else {
			f, err := os.Open(o.FilePath)
			if err != nil {
				return dasel.Value{}, fmt.Errorf("could not open file: %s: %w", o.FilePath, err)
			}
			defer f.Close()
			reader = f
		}
	}

	return storage.Load(parser, reader, options...)
}

type writeOptions struct {
	// Writer is an io.Writer that we should write to instead of FilePath.
	Writer io.Writer
	// Parser is the name of the parser we should use when reading.
	Parser string
	// FilePath is the path to the source file.
	FilePath string

	PrettyPrint bool
	Colourise   bool
	EscapeHTML  bool

	// CsvComma is the comma character used when writing CSV files.
	CsvComma string
	// CsvUseCRLF determines whether CRLF is used when writing CSV files.
	CsvUseCRLF bool

	Indent int
}

func (o *writeOptions) writeToStdout() bool {
	return o.FilePath == "" || o.FilePath == "stdout" || o.FilePath == "-"
}

func (o *writeOptions) writeParser(readOptions *readOptions) (storage.WriteParser, error) {
	if o.writeToStdout() && o.Parser == "" {
		if readOptions != nil {
			o.Parser = readOptions.Parser
		}
	}

	if o.writeToStdout() && o.Parser == "" && readOptions != nil && readOptions.FilePath != "" {
		parser, err := storage.NewWriteParserFromFilename(readOptions.FilePath)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser from read filename: %w", err)
		}
		return parser, nil
	}
	if o.Parser == "" {
		parser, err := storage.NewWriteParserFromFilename(o.FilePath)
		if err != nil {
			return nil, fmt.Errorf("could not get write parser from filename: %w", err)
		}
		return parser, nil
	}
	parser, err := storage.NewWriteParserFromString(o.Parser)
	if err != nil {
		return nil, fmt.Errorf("could not get write parser: %w", err)
	}
	return parser, nil
}

func (o *writeOptions) writeValues(cmd *cobra.Command, readOptions *readOptions, values dasel.Values) error {
	parser, err := o.writeParser(readOptions)
	if err != nil {
		return err
	}

	options := []storage.ReadWriteOption{
		storage.ColouriseOption(o.Colourise),
		storage.EscapeHTMLOption(o.EscapeHTML),
		storage.PrettyPrintOption(o.PrettyPrint),
		storage.CsvUseCRLFOption(o.CsvUseCRLF),
	}

	if o.CsvComma == "" && readOptions.CsvComma != "" {
		o.CsvComma = readOptions.CsvComma
	}

	if o.CsvComma != "" {
		options = append(options, storage.CsvCommaOption([]rune(o.CsvComma)[0]))
	}

	if o.Indent != 0 {
		options = append(options, storage.IndentOption(strconv.Itoa(o.Indent)))
	}

	writer := o.Writer
	if writer == nil {
		if o.writeToStdout() {
			writer = cmd.OutOrStdout()
		} else {
			f, err := os.Create(o.FilePath)
			if err != nil {
				return fmt.Errorf("could not open file: %s: %w", o.FilePath, err)
			}
			defer f.Close()
			writer = f
		}
	}

	for _, value := range values {
		valueBytes, err := parser.ToBytes(value, options...)
		if err != nil {
			return err
		}

		if _, err := writer.Write(valueBytes); err != nil {
			return err
		}
	}

	return nil
}

func (o *writeOptions) writeValue(cmd *cobra.Command, readOptions *readOptions, value dasel.Value) error {
	return o.writeValues(cmd, readOptions, dasel.Values{value})
}
