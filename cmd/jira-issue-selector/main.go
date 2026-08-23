package main

import (
	"fmt"
	"jira-ticket-selector/lib"
	"os"
)

func main() {
	issueId, err := lib.Run(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if issueId != "" {
		fmt.Println(issueId)
	}
}
