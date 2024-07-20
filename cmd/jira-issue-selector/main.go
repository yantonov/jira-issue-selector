package main

import (
	"fmt"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/ui"
	"os"
)

func main() {
	config := configuration.MainConfigReader{}.Load()
	err := configuration.ValidateConfig(config)
	if err != nil {
		fmt.Println(fmt.Errorf("invalid configuration: %s", err))
		os.Exit(1)
	}

	selection, err := ui.AskUser(config)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot select issue: %s", err))
		os.Exit(1)
	}

	if len(selection.TaskName) > 0 {
		fmt.Println(fmt.Sprintf("%s_%s", selection.IssueId, selection.TaskName))
	} else {
		fmt.Println(selection.IssueId)
	}

}
