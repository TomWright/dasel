package cli

import (
	"errors"
	"io"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/tomwright/dasel/v3/internal"
)

var ErrNoArgsGiven = errors.New("no arguments given")

type Globals struct {
	Stdin       io.Reader        `kong:"-"`
	Stdout      io.Writer        `kong:"-"`
	Stderr      io.Writer        `kong:"-"`
	helpPrinter kong.HelpPrinter `kong:"-"`
}

type CLI struct {
	Globals

	Query       QueryCmd       `cmd:"" default:"withargs" help:"[default] Execute a query"`
	Version     VersionCmd     `cmd:"" help:"Print the version"`
	Interactive InteractiveCmd `cmd:"" help:"Start an interactive session"`
}

func MustRun(stdin io.Reader, stdout, stderr io.Writer) {
	ctx, err := Run(stdin, stdout, stderr)
	if err == nil {
		return
	}

	if ctx == nil {
		panic(err)
	}

	ctx.Errorf("%s", err.Error())
	if errors.Is(err, ErrNoArgsGiven) {
		if err := ctx.PrintUsage(false); err != nil {
			panic(err)
		}
	}

	ctx.Exit(1)
}

func Run(stdin io.Reader, stdout, stderr io.Writer) (*kong.Context, error) {
	cli := &CLI{
		Globals: Globals{
			Stdin:       stdin,
			Stdout:      stdout,
			Stderr:      stderr,
			helpPrinter: kong.DefaultHelpPrinter,
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
		kong.TypeMapper(reflect.TypeFor[variables](), &variableMapper{}),
		kong.TypeMapper(reflect.TypeFor[extReadWriteFlags](), &extReadWriteFlagMapper{}),
		kong.OptionFunc(func(k *kong.Kong) error {
			k.Stdout = cli.Stdout
			k.Stderr = cli.Stderr
			return nil
		}),
		kong.Help(cli.helpPrinter),
	)
	err := ctx.Run()
	return ctx, err
}
