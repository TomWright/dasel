package main

import (
	"os"

	"github.com/tomwright/dasel/v3/internal/command"
)

func main() {
	cmd := command.NewRootCMD()
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}
}
