package main

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

func main() {
	config := configuration.MainConfigReader{}.Load()
	err := configuration.ValidateConfig(config)
	if err != nil {
		fmt.Println(fmt.Errorf("invalid configuration: %s", err))
		os.Exit(1)
	}

	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	selection, err := ui.AskUser(ctx, config)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println("cancelled by user")
			os.Exit(1)
		} else {
			fmt.Println(fmt.Sprintf("unexpected error: %v", err))
			os.Exit(100)
		}
	}

	if len(selection.TaskName) > 0 {
		fmt.Println(fmt.Sprintf("%s_%s", selection.IssueId, selection.TaskName))
	} else {
		fmt.Println(selection.IssueId)
	}

}
