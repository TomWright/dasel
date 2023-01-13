package main

import (
	"github.com/tomwright/dasel/v2/internal/command"
	"os"
)

func main() {
	cmd := command.NewRootCMD()
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}
}
