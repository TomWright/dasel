package cli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/tomwright/dasel/v3/parsing"
)

func NewInteractiveCmd(queryCmd *QueryCmd) *InteractiveCmd {
	return &InteractiveCmd{
		Vars:          queryCmd.Vars,
		ExtReadFlags:  queryCmd.ExtReadFlags,
		ExtWriteFlags: queryCmd.ExtWriteFlags,
		InFormat:      queryCmd.InFormat,
		OutFormat:     queryCmd.OutFormat,

		Query: queryCmd.Query,
	}
}

type InteractiveCmd struct {
	Vars          variables         `flag:"" name:"var" help:"Variables to pass to the query. E.g. --var foo=\"bar\" --var baz=json:file:./some/file.json"`
	ExtReadFlags  extReadWriteFlags `flag:"" name:"read-flag" help:"Reader flag to customise parsing. E.g. --read-flag xml-mode=structured"`
	ExtWriteFlags extReadWriteFlags `flag:"" name:"write-flag" help:"Writer flag to customise output"`
	InFormat      string            `flag:"" name:"in" short:"i" help:"The format of the input data."`
	OutFormat     string            `flag:"" name:"out" short:"o" help:"The format of the output data."`

	Query string `arg:"" help:"The query to execute." optional:"" default:""`
}

func (c *InteractiveCmd) Run(ctx *Globals) error {
	var stdInBytes []byte = nil

	if ctx.Stdin != nil {
		var err error
		stdInBytes, err = io.ReadAll(ctx.Stdin)
		if err != nil {
			return err
		}
	}

	if c.InFormat == "" && c.OutFormat == "" {
		c.InFormat = "json"
		c.OutFormat = "json"
	} else if c.InFormat == "" {
		c.InFormat = c.OutFormat
	} else if c.OutFormat == "" {
		c.OutFormat = c.InFormat
	}

	var runDasel interactiveDaselExecutor = func(selector string, root bool, formatIn parsing.Format, formatOut parsing.Format, in string) (res string, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		var stdIn *bytes.Reader = nil
		if in != "" {
			stdIn = bytes.NewReader([]byte(in))
		} else {
			stdIn = bytes.NewReader([]byte{})
		}

		o := runOpts{
			Vars:          c.Vars,
			ExtReadFlags:  c.ExtReadFlags,
			ExtWriteFlags: c.ExtWriteFlags,
			InFormat:      formatIn.String(),
			OutFormat:     formatOut.String(),
			ReturnRoot:    root,
			Unstable:      true,
			Query:         selector,

			Stdin: stdIn,
		}

		outBytes, err := run(o)
		return string(outBytes), err
	}

	p, selectorFn := newInteractiveTeaProgram(string(stdInBytes), c.Query, parsing.Format(c.InFormat), parsing.Format(c.OutFormat), runDasel)

	_, err := p.Run()
	if err != nil {
		return err
	}

	if selectorFn != nil {
		s := selectorFn()
		if s != "" {
			if _, err := fmt.Fprintf(ctx.Stdout, "%s\n", s); err != nil {
				return fmt.Errorf("error writing output: %w", err)
			}
		}
	}

	return nil
}
