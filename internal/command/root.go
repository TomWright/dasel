package command

import (
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal"
	"github.com/tomwright/dasel/internal/selfupdate"
	"os"
)

// NewRootCMD returns the root command for use with cobra.
func NewRootCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "dasel",
		Short:         "Query and modify data structures using selector strings.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.Version = internal.Version
	cmd.AddCommand(
		selectCommand(),
		putCommand(),
		deleteCommand(),
		updateCommand(selfupdate.NewUpdater(internal.Version)),
		validateCommand(),
	)
	return cmd
}

// ChangeDefaultCommand checks to see if the current os.Args target a valid subcommand.
// If they do not they are adjusted to target the given command.
// If any of the blacklisted args are set in os.Args, no action is taken.
func ChangeDefaultCommand(cmd *cobra.Command, command string, blacklistedArgs ...string) {
	subCommands := func() []string {
		results := make([]string, 0)
		for _, subCmd := range cmd.Commands() {
			results = append(results, append(subCmd.Aliases, subCmd.Name())...)
		}
		return results
	}

	if len(os.Args) > 1 {
		potentialCommand := os.Args[1]
		for _, command := range subCommands() {
			if command == potentialCommand {
				return
			}
		}
		for _, arg := range os.Args {
			for _, blacklistedArg := range blacklistedArgs {
				if arg == blacklistedArg {
					return
				}
			}
		}
		os.Args = append([]string{os.Args[0], command}, os.Args[1:]...)
	}
}
