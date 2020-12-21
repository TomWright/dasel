package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal"
	"github.com/tomwright/dasel/internal/selfupdate"
	"os"
	"runtime"
)

type updateOpts struct {
	Updater           selfupdate.Updater
	UpdateDevelopment bool
}

// ErrUnhandledGOOS is returned when there isn't a release name defined for the current OS.
var ErrUnhandledGOOS = errors.New("unhandled GOOS")

func runUpdateCommand(opts updateOpts, cmd *cobra.Command) error {
	release, err := opts.Updater.FindLatestRelease("TomWright", "dasel")
	if err != nil {
		return err
	}

	var name string
	switch runtime.GOOS {
	case "windows":
		name = "dasel_windows_amd64.exe"
	case "linux":
		name = "dasel_linux_amd64"
	case "darwin":
		name = "dasel_macos_amd64"
	default:
		return ErrUnhandledGOOS
	}

	downloadPath, err := opts.Updater.DownloadAssetWithName(release, name)
	if err != nil {
		return fmt.Errorf("could not download asset: %w", err)
	}

	currentVersion, latestVersion, err := opts.Updater.GetVersions(downloadPath)
	if err != nil {
		return fmt.Errorf("could not get version information: %w", err)
	}

	if currentVersion.IsDevelopment() && !opts.UpdateDevelopment {
		return fmt.Errorf("ignoring update for development version")
	} else {
		switch currentVersion.Compare(latestVersion) {
		case 1:
			return fmt.Errorf("your version is newer than the latest release")
		case 0:
			return fmt.Errorf("you already have the latest version")
		case -1:
			// Latest version is newer.
		}
	}

	fmt.Printf("Updating...\nCurrent version: %s\nNew version: %s\n", currentVersion, latestVersion)

	if err := opts.Updater.Replace(downloadPath); err != nil {
		_ = os.Remove(downloadPath)
		return err
	}

	fmt.Println("Successfully updated")

	return nil
}

func updateCommand() *cobra.Command {
	var updateDevFlag bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dasel to the latest stable release.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := updateOpts{
				Updater:           selfupdate.NewUpdater(internal.Version),
				UpdateDevelopment: updateDevFlag,
			}
			return runUpdateCommand(opts, cmd)
		},
	}

	cmd.Flags().BoolVar(&updateDevFlag, "dev", false, "Allow updates in development version of dasel.")

	return cmd
}
