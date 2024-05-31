package configuration

import "flag"

type CommandLineArgumentConfigLoader struct{}

func (e CommandLineArgumentConfigLoader) Load() Config {
	user := flag.String(
		"user",
		"",
		"JIRA user. Example: username@domain. Alternatively env var: "+JIRAUserEnvVar)

	hostname := flag.String(
		"hostname",
		"",
		"JIRA hostname. Example: https://company.attlassian.net. Alternatively env var: "+JIRAHostNameEnvVar)

	apikey := flag.String(
		"apikey",
		"",
		"JIRA apikey. Example: secret-key. Alternatively env var: "+JIRAApiKeyEnvVar)

	terminalStatuses := flag.String(
		"terminal-statuses",
		DefaultTerminalStatuses,
		"Terminal statuses. Alternatively env var: "+JIRATerminalStatuses)

	flag.Parse()

	return Config{
		User:             *user,
		HostName:         *hostname,
		ApiKey:           *apikey,
		TerminalStatuses: ParseTerminalStatuses(*terminalStatuses),
	}
}
