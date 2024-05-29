package main

import (
	"encoding/json"
	"fmt"
	"jira-ticket-selector/lib/jira"
	"net/http"
)

func listOfAssignedTickets(w http.ResponseWriter, req *http.Request) {
	issues := getSampleIssues()
	demo := jira.JIRAIssueListResponse{
		Total:  len(issues),
		Issues: issues,
	}
	jsonData, err := json.Marshal(demo)
	if err != nil {
		panic(err)
	}
	w.Write(jsonData)
}

func getSampleIssues() []jira.JIRAIssueResponse {
	var issues []jira.JIRAIssueResponse
	for i := 1; i <= 5; i++ {
		issues = append(issues, jira.JIRAIssueResponse{
			Key: fmt.Sprintf("PROJECT-10%d", i),
			Fields: jira.JIRAFieldsResponse{
				Summary: fmt.Sprintf("Super duper mega inspiring and important task #%d", i),
			},
		})
	}
	return issues
}

func main() {
	http.HandleFunc("/rest/api/2/search", listOfAssignedTickets)
	port := 8090
	fmt.Printf("Starting server on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
