package cli

import "github.com/tomwright/dasel/v3/internal"

type VersionCmd struct {
}

func (c *VersionCmd) Run(ctx *Globals) error {
	_, err := ctx.Stdout.Write([]byte(internal.Version + "\n"))
	return err
}
