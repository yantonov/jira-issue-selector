package main

import (
	"fmt"
	"jira-ticket-selector/lib"
	"os"
)

func main() {
	issueId, err := lib.GetIssueId()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(issueId)
}
