package cli

import (
	"fmt"
	"io"

	"github.com/tomwright/dasel/v3/execution"
	"github.com/tomwright/dasel/v3/model"
	"github.com/tomwright/dasel/v3/parsing"
)

type QueryCmd struct {
	Vars       *[]variable `flag:"" name:"var" help:"Variables to pass to the query. E.g. --var foo=\"bar\" --var baz=json:file:./some/file.json"`
	InFormat   string      `flag:"" name:"in" short:"i" help:"The format of the input data."`
	OutFormat  string      `flag:"" name:"out" short:"o" help:"The format of the output data."`
	ReturnRoot bool        `flag:"" name:"root" help:"Return the root value."`
	Unstable   bool        `flag:"" name:"unstable" help:"Allow access to potentially unstable features."`

	Query string `arg:"" help:"The query to execute." optional:"" default:""`
}

func (c *QueryCmd) Run(ctx *Globals) error {
	var opts []execution.ExecuteOptionFn

	if c.OutFormat == "" {
		c.OutFormat = c.InFormat
	}

	var reader parsing.Reader
	var err error
	if len(c.InFormat) > 0 {
		reader, err = parsing.Format(c.InFormat).NewReader()
		if err != nil {
			return fmt.Errorf("failed to get input reader: %w", err)
		}
	}

	writerOptions := parsing.DefaultWriterOptions()

	writer, err := parsing.Format(c.OutFormat).NewWriter(writerOptions)
	if err != nil {
		return fmt.Errorf("failed to get output writer: %w", err)
	}

	if c.Vars != nil {
		for _, v := range *c.Vars {
			opts = append(opts, execution.WithVariable(v.Name, v.Value))
		}
	}

	// Default to null. If stdin is being read then this will be overwritten.
	inputData := model.NewNullValue()

	var inputBytes []byte
	if ctx.Stdin != nil {
		inputBytes, err = io.ReadAll(ctx.Stdin)
		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}
	}

	if len(inputBytes) > 0 {
		if reader == nil {
			return fmt.Errorf("input format is required when reading stdin")
		}
		inputData, err = reader.Read(inputBytes)
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}
	}

	opts = append(opts, execution.WithVariable("root", inputData))

	if c.Unstable {
		opts = append(opts, execution.WithUnstable())
	}

	options := execution.NewOptions(opts...)
	out, err := execution.ExecuteSelector(c.Query, inputData, options)
	if err != nil {
		return err
	}

	if c.ReturnRoot {
		out = inputData
	}

	outputBytes, err := writer.Write(out)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	_, err = ctx.Stdout.Write(outputBytes)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}
