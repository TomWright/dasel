package cli

import (
	"io"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/internal"
)

type Globals struct {
	Stdin  io.Reader `kong:"-"`
	Stdout io.Writer `kong:"-"`
	Stderr io.Writer `kong:"-"`
}

type CLI struct {
	Globals

	Query       QueryCmd       `cmd:"" default:"withargs" help:"[default] Execute a query"`
	Version     VersionCmd     `cmd:"" help:"Print the version"`
	Interactive InteractiveCmd `cmd:"" help:"Start an interactive session"`
}

func MustRun(stdin io.Reader, stdout, stderr io.Writer) {
	ctx, err := Run(stdin, stdout, stderr)
	ctx.FatalIfErrorf(err)
}

func Run(stdin io.Reader, stdout, stderr io.Writer) (*kong.Context, error) {
	cli := &CLI{
		Globals: Globals{
			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,
		},
	}

	ctx := kong.Parse(
		cli,
		kong.Name("dasel"),
		kong.Description("Query and modify data structures from the command line."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.Vars{
			"version": internal.Version,
		},
		kong.Bind(&cli.Globals),
		kong.TypeMapper(reflect.TypeFor[*[]variable](), &variableMapper{}),
		kong.TypeMapper(reflect.TypeFor[*[]extReadWriteFlag](), &extReadWriteFlagMapper{}),
		kong.OptionFunc(func(k *kong.Kong) error {
			k.Stdout = cli.Stdout
			k.Stderr = cli.Stderr
			return nil
		}),
	)
	err := ctx.Run()
	return ctx, err
}
