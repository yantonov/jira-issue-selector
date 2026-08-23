package configuration

import (
	"flag"
	"fmt"
	"os"
)

type CommandLineArgumentConfigLoader struct{}

func (e CommandLineArgumentConfigLoader) Load(args []string) Config {
	flags := flag.NewFlagSet("jira-issue-selector", flag.ExitOnError)
	flags.Usage = func() {
		output := flags.Output()
		fmt.Fprintf(output, "Select a JIRA issue assigned to you and print its id.\n\n")

		fmt.Fprintf(output, "Usage:\n")
		PrintUsageEntries(output, []UsageEntry{
			{
				Name:        "jira-issue-selector [options]",
				Description: "select an issue, credentials are read from the keychain",
			},
			{
				Name:        "jira-issue-selector setup [setting]",
				Description: "store the credentials in the " + KeychainService + " keychain entry",
			},
		})
		fmt.Fprintln(output)

		PrintOptionGroups(flags, []OptionGroup{
			{
				Title: "Options:",
				Options: []string{
					"terminal-statuses",
					"include-ticket-title",
					"format"},
			},
		})

		fmt.Fprintf(output, "\nEnvironment variables (they override the stored settings):\n")
		PrintUsageEntries(output, EnvVarDescriptions())

		fmt.Fprintf(output, "\nExamples:\n")
		PrintUsageEntries(output, []UsageEntry{
			{
				Name:        "jira-issue-selector",
				Description: "select one of the issues assigned to you",
			},
			{
				Name:        "jira-issue-selector -format verbose",
				Description: "show the issue type in the list",
			},
			{
				Name:        "jira-issue-selector -include-ticket-title",
				Description: "use the issue title as the task name when no custom name is provided",
			},
		})

		fmt.Fprintf(output, "\nSetup:\n")
		PrintUsageEntries(output, SetupExamples())
	}

	// the options are left empty when they are not given,
	// the defaults are applied by the MainConfigReader,
	// after the environment variables have been taken into account
	terminalStatuses := flags.String(
		"terminal-statuses",
		"",
		"Terminal statuses. Default: "+DefaultTerminalStatuses)

	includeTicketTitle := flags.Bool(
		"include-ticket-title",
		false,
		"Include ticket title in the task name")

	displayFormat := flags.String(
		"format",
		"",
		"Display format for ticket list. Options: "+DisplayFormatDefault+", "+DisplayFormatVerbose+
			". Default: "+DisplayFormatDefault)

	// ExitOnError: the error is already reported by the flag set
	_ = flags.Parse(args)

	if flags.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", flags.Arg(0))
		flags.Usage()
		os.Exit(2)
	}

	config := Config{
		IncludeTicketTitle: *includeTicketTitle,
		DisplayFormat:      *displayFormat,
	}
	if *terminalStatuses != "" {
		config.TerminalStatuses = ParseTerminalStatuses(*terminalStatuses)
	}
	return config
}
