package setup

import (
	"bytes"
	"jira-ticket-selector/lib/configuration"
	"strings"
	"testing"
)

type FakeKeychain struct {
	Values map[string]string
}

func NewFakeKeychain() *FakeKeychain {
	return &FakeKeychain{Values: map[string]string{}}
}

func (e *FakeKeychain) Get(key string) (string, error) { return e.Values[key], nil }

func (e *FakeKeychain) Set(key string, value string) error {
	e.Values[key] = value
	return nil
}

func (e *FakeKeychain) Delete(key string) error {
	delete(e.Values, key)
	return nil
}

// newCommand returns a command which cannot ask anything,
// the interactive part requires a terminal.
func newCommand(keychain configuration.Keychain, output *bytes.Buffer) Command {
	return Command{
		Keychain:    keychain,
		Output:      output,
		Interactive: false,
	}
}

func TestSettingNameIsRequired(t *testing.T) {
	if _, err := settingsToAsk([]string{}); err == nil {
		t.Errorf("A missing setting name is expected to be reported, nothing is set up at once")
	}
}

func TestSingleSettingIsAskedByName(t *testing.T) {
	settings, err := settingsToAsk([]string{configuration.KeychainHostNameKey})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings) != 1 || settings[0] != configuration.KeychainHostNameKey {
		t.Errorf("Only the named setting is expected to be asked, got %v", settings)
	}
}

func TestTokenIsAnAliasForApiKey(t *testing.T) {
	settings, err := settingsToAsk([]string{TokenSettingAlias})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings) != 1 || settings[0] != configuration.KeychainApiKeyKey {
		t.Errorf("%s is expected to be an alias for %s, got %v",
			TokenSettingAlias, configuration.KeychainApiKeyKey, settings)
	}
}

func TestSettingNameIsCaseInsensitive(t *testing.T) {
	settings, err := settingsToAsk([]string{"USER"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings) != 1 || settings[0] != configuration.KeychainUserKey {
		t.Errorf("Setting name is expected to be case insensitive, got %v", settings)
	}
}

func TestUnknownSettingIsReported(t *testing.T) {
	if _, err := settingsToAsk([]string{"password"}); err == nil {
		t.Errorf("Unknown setting is expected to be reported")
	}
}

func TestOnlyOneSettingCanBeChangedAtOnce(t *testing.T) {
	_, err := settingsToAsk([]string{configuration.KeychainUserKey, configuration.KeychainHostNameKey})
	if err == nil {
		t.Errorf("Several settings at once are expected to be reported")
	}
}

func TestValueIsNotAcceptedAsAnArgument(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	err := newCommand(keychain, output).Run([]string{configuration.KeychainUserKey, "user@company.com"})
	if err == nil {
		t.Errorf("A value passed along with the setting name is expected to be reported")
	}
	if len(keychain.Values) != 0 {
		t.Errorf("Nothing is expected to be stored, got %v", keychain.Values)
	}
}

func TestNothingIsAskedWithoutTerminal(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	err := newCommand(keychain, output).Run([]string{configuration.KeychainUserKey})
	if err == nil {
		t.Errorf("Non interactive run is expected to be reported")
	}
	if len(keychain.Values) != 0 {
		t.Errorf("Nothing is expected to be stored, got %v", keychain.Values)
	}
}

func TestNothingIsSetUpWithoutASettingName(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	if err := newCommand(keychain, output).Run([]string{}); err == nil {
		t.Errorf("A setup without a setting name is expected to be reported")
	}
	if len(keychain.Values) != 0 {
		t.Errorf("Nothing is expected to be stored, got %v", keychain.Values)
	}
}

func TestTokenIsMaskedByShow(t *testing.T) {
	keychain := NewFakeKeychain()
	keychain.Values[configuration.KeychainUserKey] = "user@company.com"
	keychain.Values[configuration.KeychainApiKeyKey] = "0123456789"
	output := &bytes.Buffer{}

	if err := newCommand(keychain, output).Run([]string{"-show"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	printed := output.String()
	if strings.Contains(printed, "0123456789") {
		t.Errorf("Token is expected to be masked, got %s", printed)
	}
	if !strings.Contains(printed, "******6789") {
		t.Errorf("Last characters of the token are expected to be visible, got %s", printed)
	}
	if !strings.Contains(printed, "user@company.com") {
		t.Errorf("User is expected to be displayed, got %s", printed)
	}
	if !strings.Contains(printed, "<not set>") {
		t.Errorf("Missing hostname is expected to be displayed, got %s", printed)
	}
}

func TestSettingsAreDeleted(t *testing.T) {
	keychain := NewFakeKeychain()
	keychain.Values[configuration.KeychainUserKey] = "user@company.com"
	keychain.Values[configuration.KeychainHostNameKey] = "https://company.atlassian.net"
	keychain.Values[configuration.KeychainApiKeyKey] = "secret"
	output := &bytes.Buffer{}

	if err := newCommand(keychain, output).Run([]string{"-delete"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keychain.Values) != 0 {
		t.Errorf("Every setting is expected to be deleted, got %v", keychain.Values)
	}
}

func TestShowAndDeleteAreExclusive(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	if err := newCommand(keychain, output).Run([]string{"-show", "-delete"}); err == nil {
		t.Errorf("Conflicting options are expected to be reported")
	}
}

func TestSettingNameIsNotCombinedWithShow(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	err := newCommand(keychain, output).Run([]string{"-show", configuration.KeychainUserKey})
	if err == nil {
		t.Errorf("A setting name combined with -show is expected to be reported")
	}
}

// the flag package stops parsing the options at the first argument,
// so the option is taken as a second setting name
func TestOptionsAfterTheSettingNameAreReported(t *testing.T) {
	keychain := NewFakeKeychain()
	output := &bytes.Buffer{}

	err := newCommand(keychain, output).Run([]string{configuration.KeychainUserKey, "-show"})
	if err == nil {
		t.Errorf("An option placed after the setting name is expected to be reported")
	}
}
