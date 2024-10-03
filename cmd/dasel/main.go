package main

import (
	"os"

	"github.com/tomwright/dasel/v3/internal/cli"
)

func main() {
	cmd := cli.RootCmd()
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}
}
