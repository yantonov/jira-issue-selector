package configuration

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

	const jiraTerminalStatuses = "JIRA_TERMINAL_STATUSES"
	terminalStatuses, b := os.LookupEnv(jiraTerminalStatuses)
	if !b {
		terminalStatuses = "Done, Killed, Closed, Incomplete, Resolved"
	}
	return &Config{
		HostName:         hostname,
		User:             user,
		ApiKey:           apiKey,
		TerminalStatuses: parseTerminalStatuses(terminalStatuses),
	}, nil
}

func parseTerminalStatuses(envVar string) []string {
	tokens := strings.Split(envVar, ",")
	var result []string
	for _, token := range tokens {
		prepared := strings.TrimSpace(token)
		if strings.Contains(prepared, " ") {
			prepared = "\"" + prepared + "\""
		}
		result = append(result, prepared)
	}
	return result
}
