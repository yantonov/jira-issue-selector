package configuration

type Config struct {
	HostName         string
	User             string
	ApiKey           string
	TerminalStatuses []string
}

const DefaultTerminalStatuses = "Done, Killed, Closed, Incomplete, Resolved"
