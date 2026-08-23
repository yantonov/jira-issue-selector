package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type JIRAIssueListLoader struct{}

type JIRAIssueTypeResponse struct {
	Name string `json:"name"`
}

type JIRAFieldsResponse struct {
	Summary   string                `json:"summary"`
	IssueType JIRAIssueTypeResponse `json:"issuetype"`
}

type JIRAIssueResponse struct {
	Key    string             `json:"key"`
	Fields JIRAFieldsResponse `json:"fields"`
}

type JIRAIssueListResponse struct {
	StartAt    int                 `json:"startAt"`
	MaxResults int                 `json:"maxResults"`
	Total      int                 `json:"total"`
	Issues     []JIRAIssueResponse `json:"issues"`
}

func (e JIRAIssueListLoader) Load(config configuration.Config) (*model.IssueList, error) {
	// TODO: parameterize order by statement
	JQL := fmt.Sprintf("status not in (%s) AND assignee in (currentUser()) order by updated DESC",
		strings.Join(config.TerminalStatuses, ", "))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/rest/api/3/search/jql?jql=%s&fields=key,summary,issuetype", config.HostName, EncodeParam(JQL)),
		http.NoBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.User, config.ApiKey)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Please check your JIRA API key and ensure it's valid (see 'jira-issue-selector setup'). Invalid status code=%s response=[%s]", response.Status, responseBody)
	}
	var parsed JIRAIssueListResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return nil, err
	}
	result := ToList(parsed)
	if len(result.Issues) == 0 {
		return nil, errors.New("No tickets found. Please check your JIRA API key and ensure it's valid (see 'jira-issue-selector setup').")
	}
	return result, nil
}

func ToList(parsed JIRAIssueListResponse) *model.IssueList {
	var issues []model.Issue
	const maxSummaryLength = 80
	// TODO: add cmd param
	for _, issueItem := range parsed.Issues {
		issues = append(issues, model.Issue{
			Id:        issueItem.Key,
			Title:     trim(issueItem.Fields.Summary, maxSummaryLength),
			IssueType: issueItem.Fields.IssueType.Name,
		})
	}
	return &model.IssueList{
		Total:  parsed.Total,
		Issues: issues,
	}
}

func trim(summary string, maxLength int) string {
	if len(summary) > maxLength {
		return summary[:maxLength] + "..."
	}
	return summary
}

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}
