package configuration

import (
	"errors"
	"strings"
	"testing"
)

// clearEnvVars isolates the test from the environment it runs in.
func clearEnvVars(t *testing.T) {
	for _, name := range []string{
		JIRAUserEnvVar,
		JIRAHostNameEnvVar,
		JIRAApiKeyEnvVar,
		JIRATerminalStatuses,
		JIRAIncludeTicketTitle,
		JIRADisplayFormat} {
		t.Setenv(name, "")
	}
}

func TestCredentialsAreReadFromTheKeychain(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	keychain.Values[KeychainUserKey] = "user@company.com"
	keychain.Values[KeychainHostNameKey] = "https://company.atlassian.net"
	keychain.Values[KeychainApiKeyKey] = "secret"

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.User != "user@company.com" ||
		config.HostName != "https://company.atlassian.net" ||
		config.ApiKey != "secret" {
		t.Errorf("Credentials are expected to be read from the keychain, got %v", config)
	}
}

func TestOptionsAreTakenFromTheCommandLineArguments(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{
		"-format", DisplayFormatVerbose,
		"-include-ticket-title",
		"-terminal-statuses", "Done, Rejected"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.DisplayFormat != DisplayFormatVerbose {
		t.Errorf("Display format is expected to be taken from the arguments, got %s", config.DisplayFormat)
	}
	if !config.IncludeTicketTitle {
		t.Errorf("Ticket title is expected to be included")
	}
	if len(config.TerminalStatuses) != 2 {
		t.Errorf("Terminal statuses are expected to be parsed, got %v", config.TerminalStatuses)
	}
}

func TestDefaultOptionsAreUsed(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.DisplayFormat != DisplayFormatDefault {
		t.Errorf("Default display format is expected, got %s", config.DisplayFormat)
	}
	if len(config.TerminalStatuses) == 0 {
		t.Errorf("Default terminal statuses are expected")
	}
	if config.IncludeTicketTitle {
		t.Errorf("Ticket title is not expected to be included by default")
	}
}

func TestEnvironmentVariablesOverrideTheKeychain(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	keychain.Values[KeychainUserKey] = "stored@company.com"
	keychain.Values[KeychainHostNameKey] = "https://stored.atlassian.net"
	keychain.Values[KeychainApiKeyKey] = "stored-secret"
	t.Setenv(JIRAUserEnvVar, "env@company.com")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.User != "env@company.com" {
		t.Errorf("Environment variable is expected to win, got %s", config.User)
	}
	if config.HostName != "https://stored.atlassian.net" || config.ApiKey != "stored-secret" {
		t.Errorf("The other credentials are expected to be read from the keychain, got %v", config)
	}
}

func TestKeychainIsNotUsedWhenTheEnvironmentDefinesEveryCredential(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	keychain.Failure = errors.New("keychain is not available")
	t.Setenv(JIRAUserEnvVar, "env@company.com")
	t.Setenv(JIRAHostNameEnvVar, "https://env.atlassian.net")
	t.Setenv(JIRAApiKeyEnvVar, "env-secret")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keychain.Reads) != 0 {
		t.Errorf("Keychain is not expected to be read, got %v", keychain.Reads)
	}
	if config.User != "env@company.com" ||
		config.HostName != "https://env.atlassian.net" ||
		config.ApiKey != "env-secret" {
		t.Errorf("Credentials are expected to be read from the environment, got %v", config)
	}
}

func TestEnvironmentVariablesDefineTheOptions(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	t.Setenv(JIRADisplayFormat, DisplayFormatVerbose)
	t.Setenv(JIRAIncludeTicketTitle, "true")
	t.Setenv(JIRATerminalStatuses, "Done, Rejected")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.DisplayFormat != DisplayFormatVerbose {
		t.Errorf("Display format is expected to be taken from the environment, got %s", config.DisplayFormat)
	}
	if !config.IncludeTicketTitle {
		t.Errorf("Ticket title is expected to be included")
	}
	if len(config.TerminalStatuses) != 2 {
		t.Errorf("Terminal statuses are expected to be taken from the environment, got %v", config.TerminalStatuses)
	}
}

func TestCommandLineArgumentsOverrideTheEnvironmentVariables(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	t.Setenv(JIRADisplayFormat, DisplayFormatVerbose)
	t.Setenv(JIRATerminalStatuses, "Done, Rejected")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{
		"-format", DisplayFormatDefault,
		"-terminal-statuses", "Killed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if config.DisplayFormat != DisplayFormatDefault {
		t.Errorf("Command line argument is expected to win, got %s", config.DisplayFormat)
	}
	if len(config.TerminalStatuses) != 1 {
		t.Errorf("Command line argument is expected to win, got %v", config.TerminalStatuses)
	}
}

func TestAnUnavailableKeychainIsNotFatalWhenTheEnvironmentCompletesTheCredentials(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	keychain.Failure = errors.New("keychain is not available")
	t.Setenv(JIRAUserEnvVar, "env@company.com")
	t.Setenv(JIRAApiKeyEnvVar, "env-secret")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("An unavailable keychain is not expected to fail the load, got %v", err)
	}
	if config.User != "env@company.com" || config.ApiKey != "env-secret" {
		t.Errorf("Credentials defined by the environment are expected to be kept, got %v", config)
	}
	if config.KeychainFailure == nil {
		t.Errorf("The keychain failure is expected to be kept in the config")
	}

	err = ValidateConfig(config)
	if err == nil {
		t.Fatalf("The missing hostname is expected to be reported")
	}
	if !strings.Contains(err.Error(), JIRAHostNameEnvVar) ||
		!strings.Contains(err.Error(), "keychain is not available") {
		t.Errorf("The error is expected to report the missing credential and the keychain failure, got %v", err)
	}
}

func TestAnUnavailableKeychainIsSilentWhenEveryCredentialIsDefined(t *testing.T) {
	clearEnvVars(t)
	keychain := NewFakeKeychain()
	keychain.Failure = errors.New("keychain is not available")
	t.Setenv(JIRAUserEnvVar, "env@company.com")
	t.Setenv(JIRAHostNameEnvVar, "https://env.atlassian.net")
	t.Setenv(JIRAApiKeyEnvVar, "env-secret")

	config, err := MainConfigReader{Keychain: keychain}.Load([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateConfig(config); err != nil {
		t.Errorf("The configuration is expected to be valid, got %v", err)
	}
}

func TestMissingCredentialsAreReported(t *testing.T) {
	config := Config{TerminalStatuses: ParseTerminalStatuses(DefaultTerminalStatuses)}
	if err := ValidateConfig(config); err == nil {
		t.Errorf("Missing user is expected to be reported")
	}
}
