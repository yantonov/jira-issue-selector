package version

import (
	"runtime/debug"
	"strings"
	"testing"
)

func buildInfo(mainVersion string, settings map[string]string) func() (*debug.BuildInfo, bool) {
	return func() (*debug.BuildInfo, bool) {
		info := &debug.BuildInfo{}
		info.Main.Version = mainVersion
		for key, value := range settings {
			info.Settings = append(info.Settings, debug.BuildSetting{Key: key, Value: value})
		}
		return info, true
	}
}

func noBuildInfo() (*debug.BuildInfo, bool) {
	return nil, false
}

func TestStampedValuesWinOverTheBuildInfo(t *testing.T) {
	info := read("v1.4.2", "stamped-sha", buildInfo("v0.9.0", map[string]string{
		"vcs.revision": "build-info-sha"}))

	if info.Version != "v1.4.2" {
		t.Errorf("Stamped version is expected, got %s", info.Version)
	}
	if info.Commit != "stamped-sha" {
		t.Errorf("Stamped commit is expected, got %s", info.Commit)
	}
}

func TestCommitIsTakenFromTheBuildInfo(t *testing.T) {
	info := read("", "", buildInfo("(devel)", map[string]string{
		"vcs.revision": "abcdef",
		"vcs.modified": "false"}))

	if info.Version != DevVersion {
		t.Errorf("%s is expected for a build made from the sources, got %s", DevVersion, info.Version)
	}
	if info.Commit != "abcdef" {
		t.Errorf("Commit is expected to be taken from the build info, got %s", info.Commit)
	}
	if info.Modified {
		t.Errorf("An unmodified checkout is not expected to be reported as modified")
	}
}

func TestModifiedCheckoutIsReported(t *testing.T) {
	info := read("", "", buildInfo("(devel)", map[string]string{
		"vcs.revision": "abcdef",
		"vcs.modified": "true"}))

	if !info.Modified {
		t.Errorf("A modified checkout is expected to be reported")
	}
	if !strings.Contains(info.String(), "(modified)") {
		t.Errorf("A modified checkout is expected to be displayed, got %s", info.String())
	}
}

func TestPseudoVersionOfAnUntaggedCheckoutIsNotDisplayed(t *testing.T) {
	info := read("", "", buildInfo("v0.0.0-20260823083052-8f07524455f6+dirty", map[string]string{
		"vcs.revision": "8f07524455f6c6f63fa0fc5c873cf99de2f38796"}))

	if info.Version != DevVersion {
		t.Errorf("%s is expected instead of a pseudo version repeating the commit, got %s",
			DevVersion, info.Version)
	}
}

func TestVersionOfATaggedBuildIsKept(t *testing.T) {
	info := read("", "", buildInfo("v1.2.3", nil))

	if info.Version != "v1.2.3" {
		t.Errorf("Installed version is expected to be kept, got %s", info.Version)
	}
}

func TestUnknownCommitIsDisplayed(t *testing.T) {
	info := read("", "", noBuildInfo)

	if info.Version != DevVersion {
		t.Errorf("%s is expected without any build info, got %s", DevVersion, info.Version)
	}
	if !strings.Contains(info.String(), UnknownCommit) {
		t.Errorf("Unknown commit is expected to be displayed, got %s", info.String())
	}
}
