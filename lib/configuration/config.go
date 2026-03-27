package configuration

import "strings"

type Config struct {
	HostName           string
	User               string
	ApiKey             string
	TerminalStatuses   []string
	IncludeTicketTitle bool // if true, the ticket title will be added if no custom task name is provided
	DisplayFormat      string
}

const DefaultTerminalStatuses = "Done, Killed, Closed, Incomplete, Resolved, Canceled"
const DisplayFormatDefault = "default"
const DisplayFormatVerbose = "verbose"

func ParseTerminalStatuses(envVar string) []string {
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
