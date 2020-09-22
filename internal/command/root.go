package command

import (
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:        "dasel",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "A small helper to manage kubernetes configurations.",
}

func init() {
	RootCMD.AddCommand(
		selectCommand(),
		versionCommand(),
	)
}
