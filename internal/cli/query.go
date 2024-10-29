package cli

import "fmt"

type QueryCmd struct {
	Vars          variables         `flag:"" name:"var" help:"Variables to pass to the query. E.g. --var foo=\"bar\" --var baz=json:file:./some/file.json"`
	ExtReadFlags  extReadWriteFlags `flag:"" name:"read-flag" help:"Reader flag to customise parsing. E.g. --read-flag xml-mode=structured"`
	ExtWriteFlags extReadWriteFlags `flag:"" name:"write-flag" help:"Writer flag to customise output"`
	InFormat      string            `flag:"" name:"in" short:"i" help:"The format of the input data."`
	OutFormat     string            `flag:"" name:"out" short:"o" help:"The format of the output data."`
	ReturnRoot    bool              `flag:"" name:"root" help:"Return the root value."`
	Unstable      bool              `flag:"" name:"unstable" help:"Allow access to potentially unstable features."`
	Interactive   bool              `flag:"" name:"it" help:"Run in interactive mode."`

	Query string `arg:"" help:"The query to execute." optional:"" default:""`
}

func (c *QueryCmd) Run(ctx *Globals) error {
	if c.Interactive {
		return NewInteractiveCmd(c).Run(ctx)
	}

	o := runOpts{
		Vars:          c.Vars,
		ExtReadFlags:  c.ExtReadFlags,
		ExtWriteFlags: c.ExtWriteFlags,
		InFormat:      c.InFormat,
		OutFormat:     c.OutFormat,
		ReturnRoot:    c.ReturnRoot,
		Unstable:      c.Unstable,
		Query:         c.Query,

		Stdin: ctx.Stdin,
	}
	outBytes, err := run(o)
	if err != nil {
		return err
	}

	_, err = ctx.Stdout.Write(outBytes)
	if err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}
