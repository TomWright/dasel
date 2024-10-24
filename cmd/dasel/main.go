package main

import (
	"os"

	"github.com/tomwright/dasel/v3/internal/cli"
	_ "github.com/tomwright/dasel/v3/parsing/d"
	_ "github.com/tomwright/dasel/v3/parsing/json"
	_ "github.com/tomwright/dasel/v3/parsing/toml"
	_ "github.com/tomwright/dasel/v3/parsing/yaml"
)

func main() {
	cli.MustRun(os.Stdin, os.Stdout, os.Stderr)
}
