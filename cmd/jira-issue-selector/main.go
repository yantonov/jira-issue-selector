package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/jira"
	"os"
	"strings"
)

func SelectIssue(config configuration.Config) (string, error) {
	issues, err := jira.JIRAIssueListLoader{}.Load(config)
	if err != nil {
		return "", err
	}
	var items []string
	for _, issue := range issues.Issues {
		items = append(items, fmt.Sprintf("%s - %s", issue.Id, issue.Summary))
	}

	prompt := promptui.Select{
		Label:        "Select issue",
		Items:        items,
		HideSelected: true,
		HideHelp:     true,
		Stdout:       os.Stderr,
	}

	index, _, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return issues.Issues[index].Id, nil
}

func SelectTaskName() (string, error) {
	prompt := promptui.Prompt{
		Label:       "[Optional] task name",
		Stdout:      os.Stderr,
		HideEntered: true,
	}
	taskName, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(taskName), " ", "_"), nil
}

func main() {

	config, err := configuration.EnvVarConfigLoader{}.Load()
	if err != nil {
		fmt.Errorf("cannot load config: %s", err)
		os.Exit(1)
	}

	issueKey, err := SelectIssue(*config)
	if err != nil {
		fmt.Errorf("cannot find issueKey: %s", err)
		os.Exit(1)
	}

	issueName, err := SelectTaskName()
	if err != nil {
		fmt.Errorf("cannot select issueName %s", err)
		os.Exit(1)
	}

	if len(issueName) > 0 {
		fmt.Println(fmt.Sprintf("%s_%s", issueKey, issueName))
	} else {
		fmt.Println(issueKey)
	}

}
