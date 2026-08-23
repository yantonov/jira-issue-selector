package configuration

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type OptionGroup struct {
	Title   string
	Options []string
}

// PrintOptionGroups is needed because the flag package prints every option at once,
// in the alphabetical order.
// An option which belongs to no group is printed within the first one,
// so a newly added option is never hidden.
func PrintOptionGroups(flags *flag.FlagSet, groups []OptionGroup) {
	output := flags.Output()

	grouped := map[string]bool{}
	for _, group := range groups {
		for _, name := range group.Options {
			grouped[name] = true
		}
	}

	for i, group := range groups {
		if i > 0 {
			fmt.Fprintln(output)
		}
		if group.Title != "" {
			fmt.Fprintln(output, group.Title)
		}
		for _, name := range group.Options {
			option := flags.Lookup(name)
			if option == nil {
				continue
			}
			printOption(output, option)
		}
		if i == 0 {
			flags.VisitAll(func(option *flag.Flag) {
				if grouped[option.Name] {
					return
				}
				printOption(output, option)
			})
		}
	}
}

type UsageEntry struct {
	Name        string
	Description string
}

func PrintUsageEntries(output io.Writer, entries []UsageEntry) {
	for _, entry := range entries {
		fmt.Fprintf(output, "  %s\n", entry.Name)
		if entry.Description != "" {
			fmt.Fprintf(output, "    \t%s\n", entry.Description)
		}
	}
}

// SetupExamples is shared, the setup command is described by the usage message of both commands.
func SetupExamples() []UsageEntry {
	return []UsageEntry{
		{
			Name:        "jira-issue-selector setup",
			Description: "ask for every setting, the API token is not echoed",
		},
		{
			Name:        "jira-issue-selector setup " + KeychainApiKeyKey,
			Description: "ask for the API token only, the other settings are kept",
		},
		{
			Name:        "jira-issue-selector setup -show",
			Description: "display the stored settings, the API token is masked",
		},
		{
			Name:        "jira-issue-selector setup -delete",
			Description: "remove the stored settings",
		},
	}
}

// printOption follows the layout used by flag.PrintDefaults.
func printOption(output io.Writer, option *flag.Flag) {
	valueType, usage := flag.UnquoteUsage(option)

	var line strings.Builder
	fmt.Fprintf(&line, "  -%s", option.Name)
	if valueType != "" {
		line.WriteString(" ")
		line.WriteString(valueType)
	}
	if line.Len() <= 4 {
		line.WriteString("\t")
	} else {
		line.WriteString("\n    \t")
	}
	line.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

	if !isDefaultValueOmitted(option.DefValue) {
		if valueType == "string" {
			fmt.Fprintf(&line, " (default %q)", option.DefValue)
		} else {
			fmt.Fprintf(&line, " (default %v)", option.DefValue)
		}
	}

	fmt.Fprintln(output, line.String())
}

func isDefaultValueOmitted(defaultValue string) bool {
	switch defaultValue {
	case "", "false", "0":
		return true
	}
	return false
}
