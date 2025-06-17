package version

import (
	"fmt"
	"runtime/debug"
)

// Version information
var (
	// Version is set by build flags
	Version = "dev"
	// GitCommit is set by build flags
	GitCommit = ""
	// BuildDate is set by build flags
	BuildDate = ""
)

// Info contains version information
type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
}

// GetInfo returns version information
func GetInfo() Info {
	info := Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}

	// Try to get build info from runtime
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		info.GoVersion = buildInfo.GoVersion

		// If version is not set via build flags, try to get from VCS info
		if Version == "dev" {
			for _, setting := range buildInfo.Settings {
				switch setting.Key {
				case "vcs.revision":
					if GitCommit == "" {
						info.GitCommit = setting.Value
						if len(info.GitCommit) > 7 {
							info.GitCommit = info.GitCommit[:7]
						}
					}
				case "vcs.time":
					if BuildDate == "" {
						info.BuildDate = setting.Value
					}
				}
			}
		}
	}

	return info
}

// GetVersion returns formatted version string
func GetVersion() string {
	info := GetInfo()

	// Clean version string - remove leading 'v' if present
	version := info.Version
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}

	if info.Version == "dev" {
		if info.GitCommit != "" {
			return fmt.Sprintf("GitFleet %s (%s)", info.Version, info.GitCommit)
		}
		return fmt.Sprintf("GitFleet %s", info.Version)
	}

	return fmt.Sprintf("GitFleet v%s", version)
}

// GetVersionLong returns detailed version information
func GetVersionLong() string {
	info := GetInfo()

	// Clean version string - remove leading 'v' if present
	version := info.Version
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}

	result := fmt.Sprintf("GitFleet v%s", version)

	if info.GitCommit != "" {
		result += fmt.Sprintf(" (%s)", info.GitCommit)
	}

	if info.BuildDate != "" {
		result += fmt.Sprintf("\nBuilt: %s", info.BuildDate)
	}

	if info.GoVersion != "" {
		result += fmt.Sprintf("\nGo: %s", info.GoVersion)
	}

	return result
}
