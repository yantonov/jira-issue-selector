package lib

import (
	"context"
	"errors"
	"fmt"
	"jira-ticket-selector/lib/configuration"
	"jira-ticket-selector/lib/setup"
	"jira-ticket-selector/lib/ui"
	"jira-ticket-selector/lib/version"
	"os"
	"os/signal"

	"github.com/charmbracelet/huh"
)

const SetupCommand = "setup"
const VersionCommand = "version"

// Run returns an empty issue id when another command than the selection was executed.
func Run(args []string) (string, error) {
	keychain := configuration.SystemKeychain{}
	if len(args) > 0 {
		if args[0] == SetupCommand {
			return "", setup.Run(keychain, args[1:])
		}
		if isVersionRequested(args[0]) {
			fmt.Fprintln(os.Stdout, version.Read())
			return "", nil
		}
	}
	return GetIssueId(keychain, args)
}

// the version is requested before the options are parsed,
// so the flag package never reports the flavours of the flag as undefined
func isVersionRequested(arg string) bool {
	switch arg {
	case VersionCommand, "-v", "--v", "-version", "--version":
		return true
	}
	return false
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
