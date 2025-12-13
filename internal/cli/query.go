package cli

import "fmt"

type QueryCmd struct {
	Vars              variables         `flag:"" name:"var" help:"Variables to pass to the query. E.g. --var foo=\"bar\" --var baz=json:file:./some/file.json"`
	ExtReadWriteFlags extReadWriteFlags `flag:"" name:"rw-flag" help:"Read/Write flag to customise parsing/output. Applies to read + write E.g. --rw-flag csv-delimiter=;"`
	ExtReadFlags      extReadWriteFlags `flag:"" name:"read-flag" help:"Reader flag to customise parsing. E.g. --read-flag xml-mode=structured"`
	ExtWriteFlags     extReadWriteFlags `flag:"" name:"write-flag" help:"Writer flag to customise output. E.g. --write-flag csv-delimiter=;"`
	InFormat          string            `flag:"" name:"in" short:"i" help:"The format of the input data."`
	OutFormat         string            `flag:"" name:"out" short:"o" help:"The format of the output data."`
	ReturnRoot        bool              `flag:"" name:"root" help:"Return the root value."`
	Unstable          bool              `flag:"" name:"unstable" help:"Allow access to potentially unstable features."`
	Interactive       bool              `flag:"" name:"it" help:"Run in interactive mode."`

	ConfigPath string `name:"config" short:"c" help:"Path to config file" default:"~/dasel.yaml"`

	Query string `arg:"" help:"The query to execute." optional:"" default:""`
}

func (c *QueryCmd) Run(ctx *Globals) error {
	cfg, err := LoadConfig(c.ConfigPath)
	if err != nil {
		return err
	}

	if c.InFormat == "" && c.OutFormat == "" {
		c.InFormat = cfg.DefaultFormat
		c.OutFormat = cfg.DefaultFormat
	}

	if c.Query == "" && c.InFormat == "" && ctx.Stdin == nil {
		return ErrNoArgsGiven
	}

	if c.Interactive {
		return NewInteractiveCmd(c).Run(ctx)
	}

	o := runOpts{
		Vars:              c.Vars,
		ExtReadWriteFlags: c.ExtReadWriteFlags,
		ExtReadFlags:      c.ExtReadFlags,
		ExtWriteFlags:     c.ExtWriteFlags,
		InFormat:          c.InFormat,
		OutFormat:         c.OutFormat,
		ReturnRoot:        c.ReturnRoot,
		Unstable:          c.Unstable,
		Query:             c.Query,

		ConfigPath: c.ConfigPath,

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
