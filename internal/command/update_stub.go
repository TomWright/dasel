// +build noupdater

package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/selfupdate"
)

func updateCommand(updater *selfupdate.Updater) *cobra.Command {
	return &cobra.Command{}
}
