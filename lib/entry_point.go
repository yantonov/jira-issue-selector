package lib

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/ui"
	"os"
	"os/signal"
)

func GetIssueId() (string, error) {
	config := configuration.MainConfigReader{}.Load()
	if err := configuration.ValidateConfig(config); err != nil {
		return "", fmt.Errorf("invalid configuration: %s", err)
	}

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	selection, err := ui.AskUser(ctx, config)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", fmt.Errorf("cancelled by user")
		} else {
			return "", fmt.Errorf("unexpected error: %v", err)
		}
	}

	if len(selection.TaskName) > 0 {
		return fmt.Sprintf("%s_%s", selection.IssueId, selection.TaskName), nil
	} else {
		return fmt.Sprintf(selection.IssueId), nil
	}
}
