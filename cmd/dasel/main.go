package main

import (
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"os"
)

func main() {
	if err := command.RootCMD.Execute(); err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
}
