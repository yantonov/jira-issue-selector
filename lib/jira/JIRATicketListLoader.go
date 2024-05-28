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
	"time"
)

type JIRATicketListLoader struct{}

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

func (e JIRATicketListLoader) Load(config configuration.Config) (*model.IssueList, error) {
	const JQL = "status in (Blocked, 'In Progress', Open, Reopened, Review) AND created >= -30d AND assignee in (currentUser()) order by created DESC"
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://%s/rest/api/2/search?jql=%s", config.HostName, EncodeParam(JQL)),
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
	for _, issueItem := range parsed.Issues {
		issues = append(issues, model.Issue{
			Id:      issueItem.Key,
			Summary: issueItem.Fields.Summary,
		})
	}
	return &model.IssueList{
		Total:  parsed.Total,
		Issues: issues,
	}
}

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}
