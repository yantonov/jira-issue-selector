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

// Load leaves the undefined variables empty,
// the defaults are applied by the MainConfigReader once every source is merged.
func (e EnvVarConfigLoader) Load() Config {
	terminalStatuses := lookupEnv(JIRATerminalStatuses)

	config := Config{
		User:               lookupEnv(JIRAUserEnvVar),
		HostName:           lookupEnv(JIRAHostNameEnvVar),
		ApiKey:             lookupEnv(JIRAApiKeyEnvVar),
		IncludeTicketTitle: lookupEnv(JIRAIncludeTicketTitle) != "",
		DisplayFormat:      lookupEnv(JIRADisplayFormat),
	}
	if terminalStatuses != "" {
		config.TerminalStatuses = ParseTerminalStatuses(terminalStatuses)
	}
	return config
}

func lookupEnv(name string) string {
	value, defined := os.LookupEnv(name)
	if !defined {
		return ""
	}
	return strings.TrimSpace(value)
}

func EnvVarDescriptions() []UsageEntry {
	return []UsageEntry{
		{
			Name:        JIRAUserEnvVar,
			Description: "JIRA user, overrides the stored one",
		},
		{
			Name:        JIRAHostNameEnvVar,
			Description: "JIRA hostname, overrides the stored one",
		},
		{
			Name:        JIRAApiKeyEnvVar,
			Description: "JIRA API token, overrides the stored one",
		},
		{
			Name:        JIRATerminalStatuses,
			Description: "same as -terminal-statuses",
		},
		{
			Name:        JIRAIncludeTicketTitle,
			Description: "same as -include-ticket-title, any non empty value enables it",
		},
		{
			Name:        JIRADisplayFormat,
			Description: "same as -format",
		},
	}
}
