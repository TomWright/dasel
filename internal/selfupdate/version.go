package selfupdate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version is a semver version.
type Version struct {
	Raw   string
	Major int
	Minor int
	Patch int
}

// String returns a string representation of the Version.
func (v *Version) String() string {
	if v.IsDevelopment() || v.Major == 0 && v.Minor == 0 && v.Patch == 0 {
		return v.Raw
	}
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// IsDevelopment returns true if it's a development version.
func (v *Version) IsDevelopment() bool {
	return v.Raw == "development" || v.Raw == "dev" || strings.HasPrefix(v.Raw, "development-")
}

// Compare compares this version to the other version.
// Returns 1 if newer, -1 if older, 0 if same.
func (v *Version) Compare(other *Version) int {
	if v.Major > other.Major {
		return 1
	}
	if v.Major < other.Major {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}
	if v.Minor < other.Minor {
		return -1
	}
	if v.Patch > other.Patch {
		return 1
	}
	if v.Patch < other.Patch {
		return -1
	}

	return 0
}

var versionRegexp = regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)`)

func versionFromString(version string) *Version {
	res := &Version{
		Raw: strings.TrimSpace(version),
	}
	match := versionRegexp.FindStringSubmatch(res.Raw)
	mustAtoi := func(in string) int {
		out, err := strconv.Atoi(in)
		if err != nil {
			panic(err)
		}
		return out
	}
	if len(match) == 4 {
		res.Major = mustAtoi(match[1])
		res.Minor = mustAtoi(match[2])
		res.Patch = mustAtoi(match[3])
	}
	return res
}
