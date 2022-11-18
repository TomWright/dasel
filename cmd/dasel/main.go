package main

import (
	"github.com/tomwright/dasel/internal/command"
	"os"
)

func main() {
	cmd := command.NewRootCMD()
	// command.ChangeDefaultCommand(cmd, "select", "-v", "--version", "help")
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln("Error:", err.Error())
		os.Exit(1)
	}
}
