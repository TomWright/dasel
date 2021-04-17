// +build !noupdater

package command

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomwright/dasel/internal/selfupdate"
	"runtime"
)

type updateOpts struct {
	Updater           *selfupdate.Updater
	UpdateDevelopment bool
}

var (
	// ErrHaveLatestVersion is returned when you are already on the latest version.
	ErrHaveLatestVersion = errors.New("you already have the latest version")
	// ErrNewerVersion is returned when you are on a newer version that the latest release.
	ErrNewerVersion = errors.New("current version is newer than the latest release")
	// ErrIgnoredDev is returned when your local version is development and you have not used the --dev flag.
	ErrIgnoredDev = errors.New("ignoring update for development version")
)

func runUpdateCommand(opts updateOpts, cmd *cobra.Command) error {
	currentVersion := opts.Updater.CurrentVersion()

	out := cmd.OutOrStdout()

	_, _ = fmt.Fprintf(out, "Updating...\nCurrent version: %s\n", currentVersion)

	if currentVersion.IsDevelopment() && !opts.UpdateDevelopment {
		return ErrIgnoredDev
	}

	release, err := opts.Updater.FindLatestRelease()
	if err != nil {
		return fmt.Errorf("could not find latest release: %w", err)
	}

	releaseVersion := release.Version()
	_, _ = fmt.Fprintf(out, "Release version: %s\n", releaseVersion)

	if !currentVersion.IsDevelopment() {
		switch currentVersion.Compare(releaseVersion) {
		case 1:
			return ErrNewerVersion
		case 0:
			return ErrHaveLatestVersion
		case -1:
			// Latest version is newer.
		}
	}

	asset := release.FindAssetForSystem(runtime.GOOS, runtime.GOARCH)
	if asset == nil {
		return fmt.Errorf("could not find asset for %s %s", runtime.GOOS, runtime.GOARCH)
	}

	downloadPath, err := opts.Updater.DownloadAsset(asset)
	if err != nil {
		return fmt.Errorf("could not download asset: %w", err)
	}

	// Make sure the downloaded file is deleted if an error occurs.
	defer opts.Updater.CleanUp(downloadPath)

	latestVersion, err := opts.Updater.GetVersion(downloadPath)
	if err != nil {
		return fmt.Errorf("could not get version information: %w", err)
	}

	_, _ = fmt.Fprintf(out, "New version: %s\n", latestVersion)

	if err := opts.Updater.Replace(downloadPath); err != nil {
		return err
	}

	_, _ = fmt.Fprintln(out, "Successfully updated")

	return nil
}

func updateCommand(updater *selfupdate.Updater) *cobra.Command {
	var updateDevFlag bool
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update dasel to the latest stable release.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dev, _ := cmd.Flags().GetBool("dev")
			opts := updateOpts{
				Updater:           updater,
				UpdateDevelopment: dev,
			}
			return runUpdateCommand(opts, cmd)
		},
	}

	cmd.Flags().BoolVar(&updateDevFlag, "dev", false, "Allow updates in development version of dasel.")

	return cmd
}
