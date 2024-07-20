package ui

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/jira"
	"os"
	"strings"
)

func AskUser(config configuration.Config) (*Selection, error) {
	issues, err := jira.JIRAIssueListLoader{}.Load(config)
	if err != nil {
		return nil, err
	}

	var items []huh.Option[string]
	for _, issue := range issues.Issues {
		items = append(items, huh.NewOption(fmt.Sprintf("%s - %s", issue.Id, issue.Summary),
			issue.Id))
	}

	var selectedIssueId string
	var taskName string

	// https://github.com/charmbracelet/bubbletea/issues/860#issuecomment-2195038765
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select issue").
				Options(items...).
				Value(&selectedIssueId),
			huh.NewText().
				Title("Task name [optional]").
				CharLimit(400).
				Value(&taskName))).
		WithOutput(os.Stderr)

	formErr := form.Run()
	if formErr != nil {
		return nil, err
	}

	return &Selection{
		strings.TrimSpace(selectedIssueId),
		strings.TrimSpace(taskName),
	}, nil
}
