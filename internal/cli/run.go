package cli

import (
	"fmt"
	"io"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

type runOpts struct {
	Vars          variables
	ExtReadFlags  extReadWriteFlags
	ExtWriteFlags extReadWriteFlags
	InFormat      string
	OutFormat     string
	ReturnRoot    bool
	Unstable      bool
	Query         string

	Stdin io.Reader
}

func run(o runOpts) ([]byte, error) {
	var opts []execution.ExecuteOptionFn

	if o.OutFormat == "" && o.InFormat != "" {
		o.OutFormat = o.InFormat
	} else if o.OutFormat != "" && o.InFormat == "" {
		o.InFormat = o.OutFormat
	}

	readerOptions := parsing.DefaultReaderOptions()
	applyReaderFlags(&readerOptions, o.ExtReadFlags)

	var reader parsing.Reader
	var err error
	if len(o.InFormat) > 0 {
		reader, err = parsing.Format(o.InFormat).NewReader(readerOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to get input reader: %w", err)
		}
	}

	writerOptions := parsing.DefaultWriterOptions()
	applyWriterFlags(&writerOptions, o.ExtWriteFlags)

	writer, err := parsing.Format(o.OutFormat).NewWriter(writerOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get output writer: %w", err)
	}

	writer = parsing.MultiDocumentWriter(writer)

	opts = append(opts, variableOptions(o.Vars)...)

	// Default to null. If stdin is being read then this will be overwritten.
	inputData := model.NewNullValue()

	var inputBytes []byte
	if o.Stdin != nil {
		inputBytes, err = io.ReadAll(o.Stdin)
		if err != nil {
			return nil, fmt.Errorf("error reading stdin: %w", err)
		}
	}

	if len(inputBytes) > 0 {
		if reader == nil {
			return nil, fmt.Errorf("input format is required when reading stdin")
		}
		inputData, err = reader.Read(inputBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading input: %w", err)
		}
	}

	opts = append(opts, execution.WithVariable("root", inputData))

	if o.Unstable {
		opts = append(opts, execution.WithUnstable())
	}

	options := execution.NewOptions(opts...)
	out, err := execution.ExecuteSelector(o.Query, inputData, options)
	if err != nil {
		return nil, err
	}

	if o.ReturnRoot {
		out = inputData
	}

	outputBytes, err := writer.Write(out)
	if err != nil {
		return nil, fmt.Errorf("error writing output: %w", err)
	}

	return outputBytes, nil
}
