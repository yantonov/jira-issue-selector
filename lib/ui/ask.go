package ui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/jira"
	"jira-ticket-selector/lib/model"
	"os"
	"regexp"
	"strings"
)

func AskUser(
	ctx context.Context,
	config configuration.Config,
) (*Selection, error) {
	issues, err := jira.JIRAIssueListLoader{}.Load(config)
	if err != nil {
		return nil, err
	}

	var items []huh.Option[string]
	for _, issue := range issues.Issues {
		items = append(items, huh.NewOption(fmt.Sprintf("%s - %s", issue.Id, issue.Title),
			issue.Id))
	}

	var selectedIssueId string
	var customTaskName string

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
				Value(&customTaskName))).
		WithOutput(os.Stderr)

	if err := form.RunWithContext(ctx); err != nil {
		return nil, err
	}

	var finalTaskName = generateTaskName(issues, selectedIssueId, customTaskName, config)

	return &Selection{
		strings.TrimSpace(selectedIssueId),
		finalTaskName,
	}, nil
}

func generateTaskName(issues *model.IssueList, selectedIssueId string, customTaskName string, config configuration.Config) string {
	var normalizedTaskName = normalizeTaskName(customTaskName)
	if normalizedTaskName != "" {
		return normalizedTaskName
	}
	if config.IncludeTicketTitle {
		for _, issue := range issues.Issues {
			if issue.Id == selectedIssueId {
				return normalizeTaskName(issue.Title)
			}
		}
	}
	return ""
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func normalizeTaskName(issueTitle string) string {
	invalidCharacters := regexp.MustCompile(`[^a-zA-Z0-9_ ]+`)

	withoutSpecialChars := invalidCharacters.ReplaceAllString(issueTitle, "")
	trimmed := strings.TrimSpace(withoutSpecialChars)
	lowercased := strings.ToLower(trimmed)

	sequentialWhiteSpaces := regexp.MustCompile(` +`)
	whiteSpacesAreReplacedByUnderscore := sequentialWhiteSpaces.ReplaceAllString(lowercased, "_")
	const MaxTaskNameLength = 70
	return substr(whiteSpacesAreReplacedByUnderscore, 0, MaxTaskNameLength)
}
