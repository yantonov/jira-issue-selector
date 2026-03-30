package configuration

import (
	"fmt"
)

type MainConfigReader struct{}

func (e MainConfigReader) Load() Config {
	cmdArgConfig := CommandLineArgumentConfigLoader{}.Load()
	envVarConfig := EnvVarConfigLoader{}.Load()

	config := Config{
		User:               cmdArgConfig.User,
		HostName:           cmdArgConfig.HostName,
		ApiKey:             cmdArgConfig.ApiKey,
		TerminalStatuses:   cmdArgConfig.TerminalStatuses,
		IncludeTicketTitle: cmdArgConfig.IncludeTicketTitle,
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
	config.IncludeTicketTitle = cmdArgConfig.IncludeTicketTitle || envVarConfig.IncludeTicketTitle
	if cmdArgConfig.DisplayFormat != "" {
		config.DisplayFormat = cmdArgConfig.DisplayFormat
	} else {
		config.DisplayFormat = envVarConfig.DisplayFormat
	}
	return config
}

func ValidateConfig(config Config) error {
	if config.User == "" {
		return fmt.Errorf("User is required. You can define it using command line arg or environment variable %s", JIRAUserEnvVar)
	}
	if config.HostName == "" {
		return fmt.Errorf("Hostname is required. You can define it using command line arg or environment variable %s", JIRAHostNameEnvVar)
	}
	if config.ApiKey == "" {
		return fmt.Errorf("JIRA API KEY is required. You can define it using command line arg or environment variable %s", JIRAApiKeyEnvVar)
	}
	if len(config.TerminalStatuses) == 0 {
		return fmt.Errorf("Terminal statuses are required. You can define it using command line arg or environment variable %s", JIRATerminalStatuses)
	}
	return nil
}
