package configuration

import (
	"fmt"
)

type MainConfigReader struct {
	Keychain Keychain
}

// Load merges the settings, from the highest priority to the lowest one:
// command line arguments, environment variables, keychain, defaults.
// Credentials are not accepted as command line arguments,
// so they come either from the environment or from the keychain.
func (e MainConfigReader) Load(args []string) (Config, error) {
	config := CommandLineArgumentConfigLoader{}.Load(args)
	envVarConfig := EnvVarConfigLoader{}.Load()

	config.User = envVarConfig.User
	config.HostName = envVarConfig.HostName
	config.ApiKey = envVarConfig.ApiKey

	// a machine without a keychain daemon works when the environment defines every credential
	if !config.HasEveryCredential() {
		credentials, err := KeychainConfigLoader{Keychain: e.Keychain}.Load()
		if err != nil {
			// an unreachable keychain holds no credential, it must not discard the environment:
			// the failure is reported by ValidateConfig, and only when a credential is missing
			config.KeychainFailure = err
		}
		if config.User == "" {
			config.User = credentials.User
		}
		if config.HostName == "" {
			config.HostName = credentials.HostName
		}
		if config.ApiKey == "" {
			config.ApiKey = credentials.ApiKey
		}
	}

	if len(config.TerminalStatuses) == 0 {
		config.TerminalStatuses = envVarConfig.TerminalStatuses
	}
	if len(config.TerminalStatuses) == 0 {
		config.TerminalStatuses = ParseTerminalStatuses(DefaultTerminalStatuses)
	}

	config.IncludeTicketTitle = config.IncludeTicketTitle || envVarConfig.IncludeTicketTitle

	if config.DisplayFormat == "" {
		config.DisplayFormat = envVarConfig.DisplayFormat
	}
	if config.DisplayFormat == "" {
		config.DisplayFormat = DisplayFormatDefault
	}

	return config, nil
}

func ValidateConfig(config Config) error {
	if config.User == "" {
		return missingCredential("User", KeychainUserKey, JIRAUserEnvVar, config.KeychainFailure)
	}
	if config.HostName == "" {
		return missingCredential("Hostname", KeychainHostNameKey, JIRAHostNameEnvVar, config.KeychainFailure)
	}
	if config.ApiKey == "" {
		return missingCredential("JIRA API KEY", KeychainApiKeyKey, JIRAApiKeyEnvVar, config.KeychainFailure)
	}
	if len(config.TerminalStatuses) == 0 {
		return fmt.Errorf("Terminal statuses are required. You can define them using the -terminal-statuses command line arg")
	}
	return nil
}

func missingCredential(name string, setting string, envVar string, keychainFailure error) error {
	if keychainFailure != nil {
		return fmt.Errorf(
			"%s is required, and the stored settings could not be read: %v. Define the %s environment variable",
			name,
			keychainFailure,
			envVar)
	}
	return fmt.Errorf(
		"%s is required. Run 'jira-issue-selector setup %s' to store it in the keychain or define the %s environment variable",
		name,
		setting,
		envVar)
}
