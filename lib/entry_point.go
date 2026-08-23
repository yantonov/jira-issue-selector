package lib

import (
	"context"
	"errors"
	"fmt"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/setup"
	"jira-ticket-selector/lib/ui"
	"os"
	"os/signal"

	"github.com/charmbracelet/huh"
)

const SetupCommand = "setup"

// Run returns an empty issue id when another command than the selection was executed.
func Run(args []string) (string, error) {
	keychain := configuration.SystemKeychain{}
	if len(args) > 0 && args[0] == SetupCommand {
		return "", setup.Run(keychain, args[1:])
	}
	return GetIssueId(keychain, args)
}

func GetIssueId(keychain configuration.Keychain, args []string) (string, error) {
	config, err := configuration.MainConfigReader{Keychain: keychain}.Load(args)
	if err != nil {
		return "", err
	}
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
		return selection.IssueId, nil
	}
}
