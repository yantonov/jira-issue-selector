# JIRA issue selector

## Use case

Select JIRA issue without leaving the terminal.  
In particular, can be used to simplify creating git branches by selecting the one of assigned issues

## How to

1. put the binary to the PATH
2. define settings
```
    export JIRA_USER=username@company-domain
    export JIRA_HOSTNAME=https://company-name.attlasian.net
    export JIRA_API_KEY=private-api-key
```
3. define custom git alias, for example, like this
```
    jira="!f() { issue=$(jira-issue-selector); if [ ! $? = 0 ]; then exit 1; fi; git co -b "$issue"; }; f"
```
4. enjoy :)
```
    git jira
```

## Remarks

1. Can be used along with the commit message [hook](https://github.com/yantonov/ticket-commit-msg)
