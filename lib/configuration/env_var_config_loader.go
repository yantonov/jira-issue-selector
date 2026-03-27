package configuration

import (
	"os"
	"strings"
)

const JIRAUserEnvVar = "JIRA_USER"
const JIRAHostNameEnvVar = "JIRA_HOSTNAME"
const JIRAApiKeyEnvVar = "JIRA_API_KEY"
const JIRATerminalStatuses = "JIRA_TERMINAL_STATUSES"
const JIRAIncludeTicketTitle = "JIRA_INCLUDE_TICKET_TITLE"
const JIRADisplayFormat = "JIRA_DISPLAY_FORMAT"

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

	includeTicketTitleEnvVar, b := os.LookupEnv(JIRAIncludeTicketTitle)
	includeTicketTitle := false
	if b && strings.TrimSpace(includeTicketTitleEnvVar) != "" {
		includeTicketTitle = true
	}

	displayFormat, b := os.LookupEnv(JIRADisplayFormat)
	if !b || strings.TrimSpace(displayFormat) == "" {
		displayFormat = DisplayFormatDefault
	}

	return Config{
		HostName:           hostname,
		User:               user,
		ApiKey:             apiKey,
		TerminalStatuses:   ParseTerminalStatuses(terminalStatuses),
		IncludeTicketTitle: includeTicketTitle,
		DisplayFormat:      displayFormat,
	}
}
