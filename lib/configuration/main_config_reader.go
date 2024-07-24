package configuration

import (
	"errors"
	"fmt"
)

type MainConfigReader struct{}

func (e MainConfigReader) Load() Config {
	cmdArgConfig := CommandLineArgumentConfigLoader{}.Load()
	envVarConfig := EnvVarConfigLoader{}.Load()

	config := Config{
		User:             cmdArgConfig.User,
		HostName:         cmdArgConfig.HostName,
		ApiKey:           cmdArgConfig.ApiKey,
		TerminalStatuses: cmdArgConfig.TerminalStatuses,
	}
	if config.User == "" {
		config.User = envVarConfig.User
	}
	if config.HostName == "" {
		config.HostName = envVarConfig.HostName
	}
	if config.ApiKey == "" {
		config.ApiKey = envVarConfig.ApiKey
	}
	if len(config.TerminalStatuses) == 0 {
		config.TerminalStatuses = envVarConfig.TerminalStatuses
	}
	return config
}

func ValidateConfig(config Config) error {
	if config.User == "" {
		return errors.New(fmt.Sprintf("User is required. You can define it using command line arg or environment variable %s", JIRAUserEnvVar))
	}
	if config.HostName == "" {
		return errors.New(fmt.Sprintf("Hostname is required. You can define it using command line arg or environment variable %s", JIRAHostNameEnvVar))
	}
	if config.ApiKey == "" {
		return errors.New(fmt.Sprintf("JIRA API KEY is required. You can define it using command line arg or environment variable %s", JIRAApiKeyEnvVar))
	}
	if len(config.TerminalStatuses) == 0 {
		return errors.New(fmt.Sprintf("Terminal statuses are required. You can define it using command line arg or environment variable %s", JIRATerminalStatuses))
	}
	return nil
}
