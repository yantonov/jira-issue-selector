package lib

import "testing"

func TestVersionFlavours(t *testing.T) {
	for _, arg := range []string{"version", "-v", "--v", "-version", "--version"} {
		if !isVersionRequested(arg) {
			t.Errorf("%s is expected to request the version", arg)
		}
	}
}

func TestVersionIsNotRequested(t *testing.T) {
	for _, arg := range []string{"", "v", "setup", "-verbose", "-format", "versions"} {
		if isVersionRequested(arg) {
			t.Errorf("%s is not expected to request the version", arg)
		}
	}
}
