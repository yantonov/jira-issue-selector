# JIRA issue selector

# Table of contents

1. [Use case](#use-case)
2. [Installation](#installation)
3. [Settings](#settings)
4. [Getting JIRA API key](#getting-jira-api-key)
5. [Remarks](#remarks)

## Use case

Select JIRA issue without leaving the terminal.
In particular, can be used to simplify creating git branches by selecting the one of assigned issues.
Then you can use commit [hook](https://github.com/yantonov/ticket-commit-msg) to update commit message and include ticket number into it.

## Installation

1. Install the app to ${HOME}/bin using this script (it's assumed that ${HOME}/bin is included in PATH)

   ```bash
       curl -fsSL "https://raw.githubusercontent.com/yantonov/jira-issue-selector/master/bin/download.sh" | bash"
   ```

   or download it manually and add to the PATH
2. setup JIRA settings, they are stored in the keychain provided by your OS,
   one setting at a time

```
    jira-issue-selector setup user
    jira-issue-selector setup hostname
    jira-issue-selector setup apikey
```

3. setup git alias, for example, like this

.gitconfig

```
    [alias]
    ...
    jira="!f() { issue=$(jira-issue-selector); if [ ! $? = 0 ]; then exit 1; fi; git co -b "$issue"; }; f"
```

4. enjoy :)

```
    git jira
```

## Settings

Credentials are stored in the keychain provided by your operating system
(Keychain on macOS, Credential Manager on Windows, Secret Service on Linux),
under the `jira-issue-selector` service name, one entry per setting: `user`, `hostname`, `apikey`.

```
    jira-issue-selector setup hostname  # ask for a single setting, the other ones are kept
    jira-issue-selector setup apikey    # ask for the API token, it is not echoed
    jira-issue-selector setup -show     # display the stored settings, the API token is masked
    jira-issue-selector setup -delete   # remove the stored settings
```

The setup command requires the name of the setting to change and never accepts its value,
so a single setting is set up at a time:
the value is always asked interactively, so the API token does not end up in the shell history
and it is not echoed while being typed.
The stored value is offered as the default one, keep it by pressing enter.

| setting    | description                                            |
|------------|--------------------------------------------------------|
| `user`     | JIRA user. Example: `username@company-domain`          |
| `hostname` | JIRA hostname. Example: `https://company.attlasian.net` |
| `apikey`   | JIRA API token, `token` is accepted as an alias         |

The issue list itself is tuned by the options of the selector command:

| option                  | description                                                                                     |
|-------------------------|-------------------------------------------------------------------------------------------------|
| `-terminal-statuses`    | statuses of the issues to hide, default: `Done, Killed, Closed, Incomplete, Resolved, Canceled` |
| `-include-ticket-title` | include the ticket title in the task name when no custom name is provided                       |
| `-format`               | display format for the ticket list: `default` or `verbose`                                      |

Every setting can also be defined by an environment variable,
they override the settings stored in the keychain:

| environment variable        | overrides                                            |
|-----------------------------|------------------------------------------------------|
| `JIRA_USER`                 | the stored `user`                                    |
| `JIRA_HOSTNAME`             | the stored `hostname`                                |
| `JIRA_API_KEY`              | the stored `apikey`                                  |
| `JIRA_TERMINAL_STATUSES`    | same as `-terminal-statuses`                         |
| `JIRA_INCLUDE_TICKET_TITLE` | same as `-include-ticket-title`, any non empty value |
| `JIRA_DISPLAY_FORMAT`       | same as `-format`                                    |

Settings are merged, from the highest priority to the lowest one:
command line arguments, environment variables, keychain, defaults.
Credentials are not accepted as command line arguments,
so they come either from the environment or from the keychain.
When `JIRA_USER`, `JIRA_HOSTNAME` and `JIRA_API_KEY` are all defined,
the keychain is not accessed at all,
which is handy on machines without a keychain daemon (headless Linux, CI).

## Getting JIRA API key

For JIRA Cloud, use an Atlassian API token as the API key.

You can create or manage your API tokens here:

https://id.atlassian.com/manage-profile/security/api-tokens

Create a new API token, copy it, and use it as the API token asked by `jira-issue-selector setup apikey`.

## Remarks

1. Can be used along with the commit message [hook](https://github.com/yantonov/ticket-commit-msg)
