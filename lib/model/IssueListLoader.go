package model

import "jira-ticket-selector/lib/configuration"

type IssueListLoader interface {
	Load(config configuration.Config) (IssueList, error)
}
