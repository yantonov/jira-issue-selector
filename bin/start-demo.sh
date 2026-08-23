#!/usr/bin/env sh
set -o errexit -o nounset

cd "$(dirname "$0")/../target"

# the environment variables override the settings stored in the keychain,
# so the demo does not touch them
JIRA_USER=user JIRA_API_KEY=key JIRA_HOSTNAME=http://localhost:8090 ./jira-issue-selector
