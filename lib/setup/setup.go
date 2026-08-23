package setup

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"jira-ticket-selector/lib/configuration"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

// TokenSettingAlias is accepted as the name of the apikey setting,
// the JIRA cloud calls it an API token.
const TokenSettingAlias = "token"

type Command struct {
	Keychain    configuration.Keychain
	Output      io.Writer
	Interactive bool
}

func Run(keychain configuration.Keychain, args []string) error {
	command := Command{
		Keychain:    keychain,
		Output:      os.Stdout,
		Interactive: isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd()),
	}
	return command.Run(args)
}

func (e Command) Run(args []string) error {
	flags := flag.NewFlagSet("jira-issue-selector setup", flag.ExitOnError)
	flags.Usage = func() {
		output := flags.Output()
		fmt.Fprintf(output, "Store JIRA credentials in the keychain provided by your OS,\n")
		fmt.Fprintf(output, "under the %s service name, one entry per setting.\n", configuration.KeychainService)
		fmt.Fprintf(output, "The value is always asked interactively, it is never taken from the command line.\n\n")

		fmt.Fprintf(output, "Usage:\n")
		configuration.PrintUsageEntries(output, []configuration.UsageEntry{
			{
				Name:        "jira-issue-selector setup",
				Description: "ask for every setting",
			},
			{
				Name:        "jira-issue-selector setup <setting>",
				Description: "ask for a single setting, the other ones are kept",
			},
		})

		fmt.Fprintf(output, "\nSettings:\n")
		configuration.PrintUsageEntries(output, settingDescriptions())

		fmt.Fprintln(output)
		configuration.PrintOptionGroups(flags, []configuration.OptionGroup{
			{
				Title: "Options:",
				Options: []string{
					"show",
					"delete"},
			},
		})

		fmt.Fprintf(output, "\nExamples:\n")
		configuration.PrintUsageEntries(output, configuration.SetupExamples())
	}

	show := flags.Bool("show", false, "Show the stored settings and exit")
	deleteSettings := flags.Bool("delete", false, "Delete the stored settings and exit")

	// ExitOnError: the error is already reported by the flag set
	_ = flags.Parse(args)

	if *show && *deleteSettings {
		return fmt.Errorf("-show and -delete cannot be used together")
	}
	if *show || *deleteSettings {
		if flags.NArg() > 0 {
			return fmt.Errorf("a setting name cannot be used along with -show or -delete")
		}
		if *show {
			return e.show()
		}
		return e.delete()
	}

	settings, err := settingsToAsk(flags.Args())
	if err != nil {
		return err
	}
	return e.askAndSave(settings)
}

func settingsToAsk(args []string) ([]string, error) {
	if len(args) == 0 {
		return configuration.KeychainKeys, nil
	}
	if len(args) > 1 {
		return nil, fmt.Errorf(
			"only the name of a single setting is accepted, its value is asked interactively, got: %s",
			strings.Join(args, " "))
	}
	setting, err := settingName(args[0])
	if err != nil {
		return nil, err
	}
	return []string{setting}, nil
}

func settingName(name string) (string, error) {
	switch strings.ToLower(name) {
	case configuration.KeychainUserKey:
		return configuration.KeychainUserKey, nil
	case configuration.KeychainHostNameKey:
		return configuration.KeychainHostNameKey, nil
	case configuration.KeychainApiKeyKey, TokenSettingAlias:
		return configuration.KeychainApiKeyKey, nil
	}
	return "", fmt.Errorf("unknown setting: %s, expected one of: %s",
		name,
		strings.Join(configuration.KeychainKeys, ", "))
}

func settingDescriptions() []configuration.UsageEntry {
	return []configuration.UsageEntry{
		{
			Name:        configuration.KeychainUserKey,
			Description: "JIRA user. Example: username@domain",
		},
		{
			Name:        configuration.KeychainHostNameKey,
			Description: "JIRA hostname. Example: https://company.attlassian.net",
		},
		{
			Name:        configuration.KeychainApiKeyKey,
			Description: "JIRA API token, it is not echoed. Alias: " + TokenSettingAlias,
		},
	}
}

func (e Command) askAndSave(settings []string) error {
	if !e.Interactive {
		return fmt.Errorf("settings are asked interactively, but stdin is not a terminal")
	}

	stored, err := e.storedValues()
	if err != nil {
		return err
	}

	values, err := ask(settings, stored)
	if err != nil {
		return err
	}

	var saved []string
	for _, setting := range settings {
		value := strings.TrimSpace(values[setting])
		if value == "" || value == stored[setting] {
			continue
		}
		if err := e.Keychain.Set(setting, value); err != nil {
			return err
		}
		saved = append(saved, setting)
	}

	if len(saved) == 0 {
		fmt.Fprintf(e.Output, "Nothing changed, settings are already stored in the %s keychain entry\n",
			configuration.KeychainService)
		return nil
	}
	fmt.Fprintf(e.Output, "Saved %s to the %s keychain entry\n",
		strings.Join(saved, ", "),
		configuration.KeychainService)
	return nil
}

func ask(settings []string, stored map[string]string) (map[string]string, error) {
	var fields []huh.Field
	answers := map[string]*string{}

	for _, setting := range settings {
		answer := stored[setting]
		answers[setting] = &answer
		fields = append(fields, newInput(setting, &answer))
	}

	// keep stdout clean, the form is a part of the user interface, not of the output
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	form := huh.NewForm(huh.NewGroup(fields...)).WithOutput(os.Stderr)
	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, fmt.Errorf("cancelled by user")
		}
		return nil, err
	}

	values := map[string]string{}
	for setting, answer := range answers {
		values[setting] = *answer
	}
	return values, nil
}

func newInput(setting string, value *string) huh.Field {
	input := huh.NewInput().
		Title(title(setting)).
		Value(value).
		Validate(func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("%s is required", title(setting))
			}
			return nil
		})
	if setting == configuration.KeychainApiKeyKey {
		input = input.EchoMode(huh.EchoModePassword)
	}
	return input
}

func title(setting string) string {
	switch setting {
	case configuration.KeychainUserKey:
		return "JIRA user"
	case configuration.KeychainHostNameKey:
		return "JIRA hostname"
	case configuration.KeychainApiKeyKey:
		return "JIRA API token"
	}
	return setting
}

func (e Command) storedValues() (map[string]string, error) {
	config, err := configuration.KeychainConfigLoader{Keychain: e.Keychain}.Load()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		configuration.KeychainUserKey:     config.User,
		configuration.KeychainHostNameKey: config.HostName,
		configuration.KeychainApiKeyKey:   config.ApiKey,
	}, nil
}

func (e Command) show() error {
	stored, err := e.storedValues()
	if err != nil {
		return err
	}
	fmt.Fprintf(e.Output, "keychain service: %s\n", configuration.KeychainService)
	for _, setting := range configuration.KeychainKeys {
		value := stored[setting]
		if setting == configuration.KeychainApiKeyKey {
			value = mask(value)
		}
		fmt.Fprintf(e.Output, "%-9s %s\n", setting+":", orNotSet(value))
	}
	return nil
}

func (e Command) delete() error {
	for _, setting := range configuration.KeychainKeys {
		if err := e.Keychain.Delete(setting); err != nil {
			return err
		}
	}
	fmt.Fprintf(e.Output, "Deleted the %s keychain entry\n", configuration.KeychainService)
	return nil
}

func orNotSet(value string) string {
	if value == "" {
		return "<not set>"
	}
	return value
}

// the last characters are kept to make the token recognizable
func mask(token string) string {
	if token == "" {
		return ""
	}
	const visibleChars = 4
	asRunes := []rune(token)
	if len(asRunes) <= visibleChars {
		return strings.Repeat("*", len(asRunes))
	}
	return strings.Repeat("*", len(asRunes)-visibleChars) + string(asRunes[len(asRunes)-visibleChars:])
}
