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
	issues, err := jira.JIRATicketListLoader{}.Load(config)
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

func SelectSuffix() (string, error) {
	suffixPrompt := promptui.Prompt{
		Label:       "[Optional] suffix",
		Stdout:      os.Stderr,
		HideEntered: true,
	}
	suffix, err := suffixPrompt.Run()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(suffix), " ", "_"), nil
}

func main() {

	config, err := configuration.EnvVarConfigLoader{}.Load()
	if err != nil {
		fmt.Errorf("cannot load config: %s", err)
		os.Exit(1)
	}

	issue, err := SelectIssue(*config)
	if err != nil {
		fmt.Errorf("cannot find issue: %s", err)
		os.Exit(1)
	}

	suffix, err := SelectSuffix()
	if err != nil {
		fmt.Errorf("cannot select suffix %s", err)
		os.Exit(1)
	}

	if len(suffix) > 0 {
		fmt.Println(fmt.Sprintf("%s_%s", issue, suffix))
	} else {
		fmt.Println(issue)
	}

}
