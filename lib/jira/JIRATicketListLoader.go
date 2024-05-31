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

type JIRAFieldsResponse struct {
	Summary string `json:"summary"`
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
	JQL := fmt.Sprintf("status not in (%s) AND assignee in (currentUser()) order by created DESC",
		strings.Join(config.TerminalStatuses, ", "))
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("%s/rest/api/2/search?jql=%s", config.HostName, EncodeParam(JQL)),
		http.NoBody)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(config.User, config.ApiKey)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Invalid status code=%s response=[%s]", response.Status, responseBody))
	}
	var parsed JIRAIssueListResponse
	err = json.Unmarshal(responseBody, &parsed)
	if err != nil {
		return nil, err
	}
	return ToList(parsed), nil
}

func ToList(parsed JIRAIssueListResponse) *model.IssueList {
	var issues []model.Issue
	const maxSummaryLength = 80
	// TODO: add cmd param
	for _, issueItem := range parsed.Issues {
		issues = append(issues, model.Issue{
			Id:      issueItem.Key,
			Summary: trim(issueItem.Fields.Summary, maxSummaryLength),
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
