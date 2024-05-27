package configuration

import (
	"errors"
	"fmt"
	"os"
)

type EnvVarConfigLoader struct{}

func (e EnvVarConfigLoader) Load() (*Config, error) {
	const UserEnvVar = "JIRA_USER"
	user, b := os.LookupEnv(UserEnvVar)
	if !b {
		return nil, errors.New(fmt.Sprintf("Environment variable %s is not set", UserEnvVar))
	}

	const jiraHostName = "JIRA_HOSTNAME"
	hostname, b := os.LookupEnv(jiraHostName)
	if !b {
		return nil, errors.New(fmt.Sprintf("Environment variable %s is not set", jiraHostName))
	}

	const jiraApiKeyEnvVar = "JIRA_API_KEY"
	apiKey, b := os.LookupEnv(jiraApiKeyEnvVar)
	if !b {
		return nil, errors.New(fmt.Sprintf("Environment variable %s is not set", jiraApiKeyEnvVar))
	}

	return &Config{
		HostName: hostname,
		User:     user,
		ApiKey:   apiKey,
	}, nil
}
