package version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Version and Commit are stamped by the release build:
// go build -ldflags "-X jira-ticket-selector/lib/version.Version=v1.0.0 -X jira-ticket-selector/lib/version.Commit=<sha>"
// A build made from the sources takes them from the build info instead.
var (
	Version = ""
	Commit  = ""
)

const DevVersion = "dev"
const UnknownCommit = "<unknown>"

type Info struct {
	Version  string
	Commit   string
	Modified bool
}

func Read() Info {
	return read(Version, Commit, debug.ReadBuildInfo)
}

func read(version string, commit string, readBuildInfo func() (*debug.BuildInfo, bool)) Info {
	info := Info{
		Version: version,
		Commit:  commit,
	}

	buildInfo, available := readBuildInfo()
	if available {
		if info.Version == "" {
			info.Version = moduleVersion(buildInfo)
		}
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				if info.Commit == "" {
					info.Commit = setting.Value
				}
			case "vcs.modified":
				info.Modified = setting.Value == "true"
			}
		}
	}

	if info.Version == "" {
		info.Version = DevVersion
	}
	return info
}

// moduleVersion is a real version only when the tool is installed by 'go install module@version'.
// A build made from a checkout gets (devel), or a v0.0.0 pseudo version
// which repeats the commit reported on its own line.
func moduleVersion(buildInfo *debug.BuildInfo) string {
	version := strings.TrimSuffix(buildInfo.Main.Version, "+dirty")
	if version == "" || version == "(devel)" || strings.HasPrefix(version, "v0.0.0-") {
		return DevVersion
	}
	return version
}

func (e Info) String() string {
	commit := e.Commit
	if commit == "" {
		commit = UnknownCommit
	}
	if e.Modified {
		commit += " (modified)"
	}
	return fmt.Sprintf("version: %s\ncommit:  %s", e.Version, commit)
}
