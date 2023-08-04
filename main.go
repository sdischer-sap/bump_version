package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type VersionType string

const Major = VersionType("major")
const Minor = VersionType("minor")
const Patch = VersionType("patch")
const VERSION_FILE_LOCATION = "VERSION"

func main() {
	flag.Parse()
	args := flag.Args()
	switch len(args) {
	case 1:
		bumpType := args[0]
		handleResult(func() (*Version, error) {
			return changeVersionInFile(VersionType(bumpType))
		})
	case 2:
		bumpType := args[0]
		version := args[1]
		handleResult(func() (*Version, error) {
			return changeVersion(VersionType(bumpType), version)
		})
	default:
		os.Stdout.Write([]byte("Usage: go run <VERSION> <major|minor|patch>"))
	}

}

func handleResult(fn func() (*Version, error)) {
	v, err := fn()
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		return
	}
	os.Stdout.Write([]byte(v.String()))
}

func changeVersionInFile(vType VersionType) (*Version, error) {
	fileContent, err := os.ReadFile(VERSION_FILE_LOCATION)
	if err != nil {
		return nil, err
	}
	return changeVersion(vType, strings.TrimSpace(string(fileContent)))
}

// changeVersion takes a basic literal representing a string version, and
// increments the version number per the given VersionType.
func changeVersion(vtype VersionType, value string) (*Version, error) {
	versionNoQuotes := strings.Replace(value, "\"", "", -1)
	version, err := Parse(versionNoQuotes)
	if err != nil {
		return nil, err
	}
	if vtype == Major {
		version.Major++
		if version.Minor != -1 {
			version.Minor = 0
		}
		if version.Patch != -1 {
			version.Patch = 0
		}
	} else if vtype == Minor {
		if version.Minor == -1 {
			version.Minor = 0
		}
		if version.Patch != -1 {
			version.Patch = 0
		}
		version.Minor++
	} else if vtype == Patch {
		if version.Patch == -1 {
			version.Patch = 0
		}
		version.Patch++
	} else {
		return nil, fmt.Errorf("Invalid version type: %s", vtype)
	}
	return version, nil
}

type Version struct {
	Major int64
	Minor int64
	Patch int64
}

func (v *Version) String() string {
	if v.Major >= 0 && v.Minor >= 0 && v.Patch >= 0 {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	} else if v.Major >= 0 && v.Minor >= 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	} else if v.Major >= 0 {
		return fmt.Sprintf("%d", v.Major)
	} else {
		return "%!s(INVALID_VERSION)"
	}
}

// ParseVersion parses a version string of the forms "2", "2.3", or "0.10.11".
// Any information after the third number ("2.0.0-beta") is discarded. Very
// little effort is taken to validate the input.
//
// If a field is omitted from the string version (e.g. "0.2"), it's stored in
// the Version string as the integer -1.
func Parse(version string) (*Version, error) {
	if len(version) == 0 {
		return nil, errors.New("Empty version string")
	}

	parts := strings.SplitN(version, ".", 3)
	if len(parts) == 1 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		return &Version{
			Major: major,
			Minor: -1,
			Patch: -1,
		}, nil
	}
	if len(parts) == 2 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		minor, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		return &Version{
			Major: major,
			Minor: minor,
			Patch: -1,
		}, nil
	}
	if len(parts) == 3 {
		major, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		minor, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, err
		}
		patchParts := strings.SplitN(parts[2], "-", 2)
		patch, err := strconv.ParseInt(patchParts[0], 10, 64)
		if err != nil {
			return nil, err
		}
		return &Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}, nil
	}
	return nil, fmt.Errorf("Invalid version string: %s", version)
}
