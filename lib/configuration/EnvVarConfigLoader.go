package configuration

import (
	"os"
)

const JIRAUserEnvVar = "JIRA_USER"
const JIRAHostNameEnvVar = "JIRA_HOSTNAME"
const JIRAApiKeyEnvVar = "JIRA_API_KEY"
const JIRATerminalStatuses = "JIRA_TERMINAL_STATUSES"

type EnvVarConfigLoader struct{}

func (e EnvVarConfigLoader) Load() Config {
	user, b := os.LookupEnv(JIRAUserEnvVar)
	if !b {
		user = ""
	}

	hostname, b := os.LookupEnv(JIRAHostNameEnvVar)
	if !b {
		hostname = ""
	}

	apiKey, b := os.LookupEnv(JIRAApiKeyEnvVar)
	if !b {
		apiKey = ""
	}

	terminalStatuses, b := os.LookupEnv(JIRATerminalStatuses)
	if !b {
		terminalStatuses = DefaultTerminalStatuses
	}
	return Config{
		HostName:         hostname,
		User:             user,
		ApiKey:           apiKey,
		TerminalStatuses: ParseTerminalStatuses(terminalStatuses),
	}
}
