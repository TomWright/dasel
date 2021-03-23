package internal

import (
	"runtime/debug"
)

// Version represents the current version of dasel.
// The real version number is injected at build time using ldflags.
var Version = "development"

func init() {
	// Version is set by ldflags on build.
	if Version != "development" {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	// https://github.com/golang/go/issues/29228
	if info.Main.Version == "(devel)" || info.Main.Version == "" {
		return
	}

	Version += "-" + info.Main.Version
}
