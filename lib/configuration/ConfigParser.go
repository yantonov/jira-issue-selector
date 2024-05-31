package configuration

import "strings"

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
