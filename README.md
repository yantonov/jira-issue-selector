# JIRA issue selector

# Table of contents
1. [Use case](#use-case)
2. [Installation](#installation)
3. [Remarks](#remarks)


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
2. setup JIRA settings
```
    export JIRA_USER=username@company-domain
    export JIRA_HOSTNAME=https://company-name.attlasian.net
    export JIRA_API_KEY=private-api-key
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

## Remarks

1. Can be used along with the commit message [hook](https://github.com/yantonov/ticket-commit-msg)
